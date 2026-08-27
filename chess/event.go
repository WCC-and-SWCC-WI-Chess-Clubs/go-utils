package chess

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Event represents a USCF-rated chess event
type Event struct {
	ID           string
	Href         string
	StartDate    time.Time
	FinishDate   time.Time
	name         string
	nameOverride string
	City         string
	State        string
	Players      string
	Sections     string
}

func NewEvent() *Event {
	return &Event{
		StartDate:  time.Now(),
		FinishDate: time.Now(),
	}
}

func (e *Event) GetName() string {
	if len(e.nameOverride) > 0 {
		return e.nameOverride
	}
	return e.name
}

func (e *Event) SetName(aName string) {
	if strings.TrimSpace(aName) != "" {
		// always coerce the Tournament Name to title case
		e.name = toTitleCase(aName)
	} else {
		err := errors.New("Tournament name may not be empty")
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func (e *Event) SetNameOverride(aName string) {
	if len(aName) > 0 {
		e.nameOverride = aName
		if aName == toTitleCase(e.name) {
			fmt.Printf("Name override not necessary for Event (%s) %s \n", e.ID, aName)
		}
	}
}

func (e *Event) String() string {
	return "id: " + e.ID + "\n" +
		"href: " + e.Href + "\n" +
		"start: " + e.StartDate.Format("2006-01-02") + "\n" +
		"finish: " + e.FinishDate.Format("2006-01-02") + "\n" +
		"name: " + e.name + "\n" +
		"sections: " + e.Sections + "\n" +
		"players: " + e.Players + "\n" +
		"city: " + e.City + "\n" +
		"state: " + e.State + "\n\n"
}

func (e *Event) SetHref(href string) {
	e.Href = href
	if len(href) > 7 {
		e.ID = href[7:]
	}
}

func (e *Event) SetStartDate(dateStr string) error {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}
	e.StartDate = t
	return nil
}

func (e *Event) SetFinishDate(dateStr string) error {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}
	e.FinishDate = t
	return nil
}

func toTitleCase(s string) string {
	// Create an English title caser
	caser := cases.Title(language.English)
	output := caser.String(s)
	return output
}
