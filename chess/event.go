package chess

import "time"

// Event represents a USCF-rated chess event
type Event struct {
	ID           string
	Href         string
	StartDate    time.Time
	FinishDate   time.Time
	Name         string
	NameOverride string
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
	if len(e.NameOverride) > 0 {
		return e.NameOverride
	}
	return e.Name
}

func (e *Event) String() string {
	return "id: " + e.ID + "\n" +
		"href: " + e.Href + "\n" +
		"start: " + e.StartDate.Format("2006-01-02") + "\n" +
		"finish: " + e.FinishDate.Format("2006-01-02") + "\n" +
		"name: " + e.Name + "\n" +
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
