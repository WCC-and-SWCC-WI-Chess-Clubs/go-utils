package chess

// Tournament holds data fetched from the USCF Ratings API for a single event
type Tournament struct {
	Event    EventResponse
	Sections []SectionResponse
	Results  []StandingsResponse
}

func NewTournament() *Tournament {
	return &Tournament{}
}

// --- API response types (ratings-api.uschess.org) ---

// EventResponse is the JSON shape from GET /api/v1/rated-events/{id}
type EventResponse struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	StartDate    string       `json:"startDate"`
	EndDate      string       `json:"endDate"`
	SectionCount int          `json:"sectionCount"`
	PlayerCount  int          `json:"playerCount"`
	GameCount    int          `json:"gameCount"`
	StateCode    string       `json:"stateCode"`
	City         string       `json:"city"`
	ZipCode      string       `json:"zipCode"`
	Status       string       `json:"status"`
	Sections     []SectionRef `json:"sections"`
}

// SectionRef is the lightweight section summary embedded in EventResponse
type SectionRef struct {
	ID     string `json:"id"`
	Number int    `json:"number"`
	Name   string `json:"name"`
}

// SectionResponse is the JSON shape from GET /api/v1/rated-events/{id}/sections/{n}
type SectionResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Number       int    `json:"number"`
	PlayerCount  int    `json:"playerCount"`
	GameCount    int    `json:"gameCount"`
	RoundCount   int    `json:"roundCount"`
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	Format       string `json:"format"`
	TimeControl  string `json:"timeControl"`
	RatingSystem string `json:"ratingSystem"`
	IsGrandPrix  bool   `json:"isGrandPrix"`
	IsBlitz      bool   `json:"isBlitz"`
}

// StandingsResponse is the JSON shape from GET /api/v1/rated-events/{id}/sections/{n}/standings
type StandingsResponse struct {
	Items []PlayerResult `json:"items"`
}

// PlayerResult is one player's standing within a section
type PlayerResult struct {
	Ordinal       int            `json:"ordinal"`
	PairingNumber int            `json:"pairingNumber"`
	PlayerID      string         `json:"playerId"`
	MemberID      string         `json:"memberId"`
	FirstName     string         `json:"firstName"`
	LastName      string         `json:"lastName"`
	StateRep      string         `json:"stateRep"`
	Score         float64        `json:"score"`
	RoundOutcomes []RoundOutcome `json:"roundOutcomes"`
	Ratings       []RatingEntry  `json:"ratings"`
}

// RoundOutcome describes one round's result for a player
type RoundOutcome struct {
	ID                string `json:"id"`
	RoundNumber       int    `json:"roundNumber"`
	Outcome           string `json:"outcome"`
	Color             string `json:"color"`
	OpponentOrdinal   int    `json:"opponentOrdinal"`
	OpponentMemberID  string `json:"opponentMemberId"`
	OpponentFirstName string `json:"opponentFirstName"`
	OpponentLastName  string `json:"opponentLastName"`
}

// RatingEntry holds pre/post ratings for one rating system
type RatingEntry struct {
	PreRating    int    `json:"preRating"`
	PostRating   int    `json:"postRating"`
	RatingSystem string `json:"ratingSystem"`
}
