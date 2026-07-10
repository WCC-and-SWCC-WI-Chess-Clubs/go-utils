package chess

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"chess-utils/utils"
)

const (
	muirBaseCrosstableURL = "https://ratings.uschess.org/event/"
	SwccAffiliateID       = "A6011047"
	WccAffiliateID        = "A5008948"
)

// Reader provides high-level access to USCF event data
type Reader struct{}

func NewReader() *Reader { return &Reader{} }

// GetCrosstableURL returns the public crosstable URL for an event
func (r *Reader) GetCrosstableURL(eventID string) string {
	return muirBaseCrosstableURL + eventID
}

// GetPlayerName formats a player's name from a PlayerResult as display HTML
func (r *Reader) GetPlayerName(p PlayerResult) string {
	name := utils.MassageName(p.FirstName + " " + p.LastName)
	return strings.ReplaceAll(name, " ", "&nbsp;")
}

// GetEventWinnersNames returns a space-joined string of winner names for the event's first section
func (r *Reader) GetEventWinnersNames(eventID string) (string, error) {
	api := NewRatingsAPI()
	tournament, err := api.QueryEvent(eventID)
	if err != nil {
		return "", err
	}
	if len(tournament.Results) == 0 {
		return "", fmt.Errorf("no results for event %s", eventID)
	}
	results := tournament.Results[0]

	winScore := -1.0
	outName := ""
	for _, player := range results.Items {
		if player.Score >= winScore {
			winScore = player.Score
			if outName != "" {
				outName += " "
			}
			outName += r.GetPlayerName(player)
		} else {
			break
		}
	}
	return outName, nil
}

// WinnersEntry holds one year's champion info for both the club champ and memorial
type WinnersEntry struct {
	Year      int
	MemLink   string
	MemWinner string
	CCLink    string
	CCWinner  string
}

type winnersJSON struct {
	Memorial []winnerRef `json:"memorial"`
	ClubChp  []winnerRef `json:"clubChp"`
}

type winnerRef struct {
	EventYear string `json:"eventYear"`
	EventID   string `json:"eventID"`
	Winner    string `json:"winner"`
}

// GetWinners loads winners.json and fetches names from the API for events with real IDs.
// Note: this can take ~10 minutes due to API rate limiting.
func (r *Reader) GetWinners() ([]WinnersEntry, error) {
	data, err := os.ReadFile("data/Winners.json")
	if err != nil {
		return nil, fmt.Errorf("read Winners.json: %w", err)
	}
	var winners winnersJSON
	if err := json.Unmarshal(data, &winners); err != nil {
		return nil, err
	}

	var outSet []WinnersEntry
	currentYear := time.Now().Year()

	for year := currentYear; year > 1950; year-- {
		yearStr := fmt.Sprintf("%d", year)
		var memRef, ccRef *winnerRef
		for i := range winners.Memorial {
			if winners.Memorial[i].EventYear == yearStr {
				memRef = &winners.Memorial[i]
				break
			}
		}
		for i := range winners.ClubChp {
			if winners.ClubChp[i].EventYear == yearStr {
				ccRef = &winners.ClubChp[i]
				break
			}
		}
		if memRef == nil && ccRef == nil {
			return outSet, nil
		}

		entry := WinnersEntry{Year: year}

		if memRef != nil {
			if memRef.EventID == "x" {
				entry.MemWinner = memRef.Winner
			} else {
				name, err := r.GetEventWinnersNames(memRef.EventID)
				if err != nil {
					return nil, err
				}
				entry.MemWinner = name
				entry.MemLink = r.GetCrosstableURL(memRef.EventID)
			}
		}

		if ccRef != nil {
			if ccRef.EventID == "x" {
				entry.CCWinner = ccRef.Winner
			} else {
				name, err := r.GetEventWinnersNames(ccRef.EventID)
				if err != nil {
					return nil, err
				}
				entry.CCWinner = name
				entry.CCLink = r.GetCrosstableURL(ccRef.EventID)
			}
		}

		outSet = append(outSet, entry)
		time.Sleep(7500 * time.Millisecond)
	}
	return outSet, nil
}

// GetPastEvent fetches a single tournament by ID from the ratings API
func (r *Reader) GetPastEvent(tournamentID string) (*Tournament, error) {
	api := NewRatingsAPI()
	return api.QueryEvent(tournamentID)
}

// GetPastEvents fetches all WCC-affiliated events from the ratings API
func (r *Reader) GetPastEvents(theAffiliateID string) ([]*Event, error) {
	api := NewRatingsAPI()
	return api.QueryEvents(theAffiliateID)
}

// GetPastEventsForAffiliate fetches all events for the given affiliate ID from the ratings API
func (r *Reader) GetPastEventsForAffiliate(affiliateID string) ([]*Event, error) {
	api := NewRatingsAPI()
	return api.QueryEventsByAffiliate(affiliateID)
}
