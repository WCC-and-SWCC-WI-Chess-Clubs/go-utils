package chess

import "strings"

// Section holds a tournament section's metadata and player list
type Section struct {
	Name     string
	UscfName string
	Href     string
	Players  []*Player
}

func NewSection(name, uscfName, href string) *Section {
	return &Section{Name: name, UscfName: uscfName, Href: href}
}

func (s *Section) AddPlayer(p *Player) {
	s.Players = append(s.Players, p)
}

func (s *Section) GetPlayerCount() int {
	return len(s.Players)
}

func (s *Section) GetRoundCount() int {
	if len(s.Players) > 0 {
		return len(s.Players[0].Rounds)
	}
	return 0
}

func (s *Section) ToHTML() string {
	var b strings.Builder
	for _, p := range s.Players {
		b.WriteString(p.ToHTML())
	}
	return b.String()
}

func (s *Section) GetNameHTML() string {
	return "<span class='wccSection'>" + s.Name + "</span>\n"
}
