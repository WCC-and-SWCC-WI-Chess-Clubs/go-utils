package chess

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"time"
)

const (
	muirAPIBaseURL = "https://ratings-api.uschess.org/api/v1"
	muirEventURL   = "https://ratings.uschess.org/event/"
	//swccAffiliateID = "A6011047"
	//wccAffiliateID  = "A5008948"
	apiPageSize = 100
)

// RatingsAPI is a client for the USCF Ratings API
type RatingsAPI struct {
	SleepSeconds int
}

func NewRatingsAPI() *RatingsAPI {
	return &RatingsAPI{SleepSeconds: 4}
}

func (r *RatingsAPI) fetchJSON(url string, target any) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("GET %s: %w", url, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}

// affiliateEventsPage is the paginated response from /affiliates/{id}/events
type affiliateEventsPage struct {
	Items []affiliateEventItem `json:"items"`
}

type affiliateEventItem struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	SectionCount int    `json:"sectionCount"`
	PlayerCount  int    `json:"playerCount"`
	StateCode    string `json:"stateCode"`
	City         string `json:"city"`
}

func (r *RatingsAPI) createEvent(item affiliateEventItem) *Event {
	e := NewEvent()
	e.ID = item.ID
	e.Name = item.Name
	e.City = item.City
	e.State = item.StateCode
	e.Sections = strconv.Itoa(item.SectionCount)
	e.Players = strconv.Itoa(item.PlayerCount)
	e.Href = muirEventURL + e.ID
	_ = e.SetStartDate(item.StartDate)
	_ = e.SetFinishDate(item.EndDate)
	return e
}

func (r *RatingsAPI) buildUrl(theAffiliateID string, theOffset int, theCutoffYear int) string {
	var url string
	if theCutoffYear >= 1990 {
		fromDate := fmt.Sprintf("%d-01-01", theCutoffYear)
		url = fmt.Sprintf("%s/affiliates/%s/events?FromDate=%s&SortBy=StartDate&Offset=%d&Size=%d",
			muirAPIBaseURL, theAffiliateID, fromDate, theOffset, apiPageSize)
	} else {
		url = fmt.Sprintf("%s/affiliates/%s/events?SortBy=StartDate&Offset=%d&Size=%d",
			muirAPIBaseURL, theAffiliateID, theOffset, apiPageSize)
	}
	return url
}

// QueryEvents retrieves all WCC-affiliated rated events, sorted newest-first
func (r *RatingsAPI) QueryEvents(theAffiliateID string, theCutoffYear int) ([]*Event, error) {
	var eventList []*Event
	for index := 0; ; index++ {
		if index > 0 {
			fmt.Println("sleeping...")
			time.Sleep(time.Duration(r.SleepSeconds) * time.Second)
		}
		fmt.Println(time.Now())
		offset := index * apiPageSize
		url := r.buildUrl(theAffiliateID, offset, theCutoffYear)

		var page affiliateEventsPage
		if err := r.fetchJSON(url, &page); err != nil {
			return nil, err
		}
		if len(page.Items) == 0 {
			break
		}
		for _, item := range page.Items {
			eventList = append(eventList, r.createEvent(item))
		}
		fmt.Printf("Number of events: %d\n", len(eventList))
	}

	sort.Slice(eventList, func(i, j int) bool {
		return eventList[i].FinishDate.After(eventList[j].FinishDate)
	})
	return eventList, nil
}

// QueryEvent fetches a tournament by ID including all sections and standings
func (r *RatingsAPI) QueryEvent(tournamentID string) (*Tournament, error) {
	url := fmt.Sprintf("%s/rated-events/%s", muirAPIBaseURL, tournamentID)
	var eventResp EventResponse
	if err := r.fetchJSON(url, &eventResp); err != nil {
		return nil, fmt.Errorf("fetch event: %w", err)
	}

	t := NewTournament()
	t.Event = eventResp

	for i := range eventResp.Sections {
		sectionNum := strconv.Itoa(i + 1)
		sec, err := r.FetchSection(tournamentID, sectionNum)
		if err != nil {
			return nil, fmt.Errorf("fetch section %s: %w", sectionNum, err)
		}
		t.Sections = append(t.Sections, sec)

		results, err := r.FetchSectionResults(tournamentID, sectionNum)
		if err != nil {
			return nil, fmt.Errorf("fetch results %s: %w", sectionNum, err)
		}
		t.Results = append(t.Results, results)
	}
	return t, nil
}

// FetchSection fetches metadata for one section of an event
func (r *RatingsAPI) FetchSection(tournamentID, sectionID string) (SectionResponse, error) {
	url := fmt.Sprintf("%s/rated-events/%s/sections/%s", muirAPIBaseURL, tournamentID, sectionID)
	var sec SectionResponse
	err := r.fetchJSON(url, &sec)
	return sec, err
}

// FetchSectionResults fetches the standings for one section of an event
func (r *RatingsAPI) FetchSectionResults(tournamentID, sectionID string) (StandingsResponse, error) {
	url := fmt.Sprintf("%s/rated-events/%s/sections/%s/standings", muirAPIBaseURL, tournamentID, sectionID)
	var results StandingsResponse
	err := r.fetchJSON(url, &results)
	return results, err
}
