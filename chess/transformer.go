package chess

import (
	"fmt"
	"io"
)

// Transformer converts Tournament API data into blog-post HTML
type Transformer struct{}

func NewTransformer() *Transformer { return &Transformer{} }

// CreateBlogResults prints HTML for all sections of a tournament to w
func (t *Transformer) CreateBlogResults(w io.Writer, tournament *Tournament) {
	for i := range tournament.Sections {
		t.printSectionResults(w, tournament.Sections[i], tournament.Results[i])
		t.printSectionDivider(w)
	}
}

func (t *Transformer) printSectionDivider(w io.Writer) {
	fmt.Fprintln(w, "<div><br/></div>")
}

func (t *Transformer) determineRounds(results StandingsResponse) int {
	if len(results.Items) > 0 {
		return len(results.Items[0].RoundOutcomes)
	}
	return 0
}

func (t *Transformer) getPlayerName(p PlayerResult) string {
	return p.FirstName + "&nbsp;" + p.LastName
}

// getRatingByID finds a rating entry by ratingSystem ("R" = Regular, "Q" = Quick)
func (t *Transformer) getRatingByID(ratings []RatingEntry, ratingID string) *RatingEntry {
	for i := range ratings {
		if ratings[i].RatingSystem == ratingID {
			return &ratings[i]
		}
	}
	if len(ratings) > 0 {
		return &ratings[0]
	}
	return nil
}

func (t *Transformer) printRound(w io.Writer, round RoundOutcome) {
	abbrevs := map[string]string{
		"ByeFull":    "B",
		"ByeHalf":    "H",
		"Draw":       "D",
		"Forfeit":    "F",
		"Loss":       "L",
		"Unpaired":   "U",
		"Win":        "W",
		"WinForfeit": "X",
	}
	result := abbrevs[round.Outcome]
	if round.OpponentOrdinal > 0 {
		result += fmt.Sprintf("%d", round.OpponentOrdinal)
	}
	fmt.Fprintf(w, "<td>%s</td>\n", result)
}

func (t *Transformer) printSectionResults(w io.Writer, section SectionResponse, results StandingsResponse) {
	rounds := t.determineRounds(results)
	fmt.Fprintln(w, "<div class='wccSectionResults'>")
	fmt.Fprintf(w, "<span class='wccSection'>%s</span>\n", section.Name)
	fmt.Fprintln(w, "<table class='wccCrosstable'>")
	fmt.Fprintln(w, "<thead>")
	fmt.Fprintln(w, "<tr>")
	fmt.Fprintln(w, "<th>No.</th>")
	fmt.Fprintln(w, "<th>Player Name</th>")
	fmt.Fprintln(w, "<th>State</th>")
	fmt.Fprintln(w, "<th>USCF ID</th>")
	fmt.Fprintln(w, "<th>Pre</th>")
	fmt.Fprintln(w, "<th>Post</th>")
	for i := 1; i <= rounds; i++ {
		fmt.Fprintf(w, "<th>R%d</th>\n", i)
	}
	fmt.Fprintln(w, "<th>Total</th>")
	fmt.Fprintln(w, "</tr>")
	fmt.Fprintln(w, "</thead>")
	fmt.Fprintln(w, "<tbody>")
	for _, player := range results.Items {
		t.printPlayer(w, player)
	}
	fmt.Fprintln(w, "</tbody>")
	fmt.Fprintln(w, "</table>")
	fmt.Fprintln(w, "</div>")
}

func (t *Transformer) printPlayer(w io.Writer, player PlayerResult) {
	fmt.Fprintln(w, "<tr>")
	fmt.Fprintf(w, "<td>%d</td>\n", player.Ordinal)
	fmt.Fprintf(w, "<td>%s</td>\n", t.getPlayerName(player))
	fmt.Fprintf(w, "<td>%s</td>\n", player.StateRep)
	fmt.Fprintf(w, "<td>%s</td>\n", player.MemberID)
	t.printPlayerRatings(w, player)
	for _, round := range player.RoundOutcomes {
		t.printRound(w, round)
	}
	fmt.Fprintf(w, "<td>%g</td>\n", player.Score)
	fmt.Fprintln(w, "</tr>")
}

func (t *Transformer) printPlayerRatings(w io.Writer, player PlayerResult) {
	rating := t.getRatingByID(player.Ratings, "R")
	if rating == nil {
		fmt.Fprintln(w, "<td>&nbsp;</td><td>&nbsp;</td>")
		return
	}
	if rating.PreRating > 0 {
		fmt.Fprintf(w, "<td>%d</td>\n", rating.PreRating)
	} else {
		fmt.Fprintln(w, "<td>Unr</td>")
	}
	fmt.Fprintf(w, "<td>%d</td>\n", rating.PostRating)
}
