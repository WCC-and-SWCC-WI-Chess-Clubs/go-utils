package chess

import (
	"fmt"
	"strings"

	"chess-utils/utils"
)

var normDisp = map[string]string{
	"N:4":  "1200",
	"N:3":  "1400",
	"N:2":  "1600",
	"N:1":  "1800",
	"N:CM": "2000",
	"N:LM": "2200",
	"N:SM": "2400",
	"N:C":  "2000",
	"N:M":  "2200",
	"N:S":  "2400",
	"":     "",
}

// Player represents a tournament participant
type Player struct {
	ID       string
	Rank     string
	Name     string
	State    string
	RatePre  string
	RatePost string
	Rounds   []string
	RdScore  []float64
	RdTotal  []float64
	Norm     string
	Total    string
}

func NewPlayer() *Player {
	return &Player{}
}

// CreatePlayerFromWinTDXtbl parses a WinTD crosstable line, e.g.:
//
//	   1.    Waller, Matt (1) .............  WI     2058 W10   W6    W2    -N-     3.0
func CreatePlayerFromWinTDXtbl(line string, numRounds int) *Player {
	p := NewPlayer()
	if line == "" {
		return p
	}
	// Pad so index math doesn't panic
	for len(line) < 53+(6*(numRounds+1)) {
		line += " "
	}

	p.Name = strings.Trim(line[9:38], ". ")
	p.RatePre = strings.TrimSpace(line[48:52])

	lastGood := -1
	for i := 0; i <= numRounds; i++ {
		start := 53 + (6 * i)
		end := start + 6
		if end > len(line) {
			end = len(line)
		}
		thisTag := ""
		if start < len(line) {
			thisTag = strings.TrimSpace(line[start:end])
		}
		if thisTag != "" {
			lastGood = i
		}
		p.Rounds = append(p.Rounds, thisTag)
	}
	// If we didn't find the total where expected (e.g. fewer rounds played), clear the stray entry
	if lastGood != numRounds && lastGood >= 0 {
		p.Rounds[lastGood] = ""
	}

	p.CalculateRoundScores()
	return p
}

// CreatePlayer parses a player from two USCF HTML crosstable lines
func CreatePlayer(line1, line2 string) *Player {
	p := NewPlayer()
	parseLineOne(p, line1)
	parseLineTwo(p, line2)
	return p
}

func parseLineOne(p *Player, line string) {
	line = strings.TrimSpace(line)
	items := strings.Split(line, "|")
	if len(items) >= 2 {
		p.Rank = utils.GetText(items[0])
		p.Name = utils.MassageName(utils.GetText(items[1]))
	}
	if len(items) > 2 {
		p.Total = items[2]
		for i := 3; i < len(items)-1; i++ {
			p.Rounds = append(p.Rounds, items[i])
		}
	}
}

func parseLineTwo(p *Player, line string) {
	line = strings.TrimSpace(line)
	if len(line) < 2 {
		return
	}
	p.State = line[0:2]
	if len(line) > 13 {
		p.ID = strings.TrimSpace(line[5:13])
	}
	rating := ""
	if len(line) > 35 {
		rating = line[19:35]
	}
	if !strings.Contains(rating, "->") {
		p.RatePre = "0"
		p.RatePost = "0"
	} else {
		pos := strings.Index(rating, "->")
		p.RatePre = strings.TrimSpace(rating[:pos])
		p.RatePost = strings.TrimSpace(rating[pos+2:])
	}
	items := strings.Split(line, "|")
	if len(items) > 2 {
		key := strings.TrimSpace(items[2])
		if val, ok := normDisp[key]; ok {
			p.Norm = val
		}
	}
	for i := range p.Rounds {
		idx := 3 + i
		if idx < len(items) && strings.TrimSpace(items[idx]) != "" {
			p.Rounds[i] = p.Rounds[i] + "/" + strings.TrimSpace(items[idx])
		}
	}
}

// CalculateRoundScores computes running score totals from round results
func (p *Player) CalculateRoundScores() {
	p.RdScore = nil
	p.RdTotal = nil
	var curTotal float64
	for _, rd := range p.Rounds {
		if len(rd) > 0 && (rd[0] == 'W' || rd == "BYE" || rd == "-B-" || rd[0] == 'X') {
			p.RdScore = append(p.RdScore, 1.0)
			curTotal += 1.0
		} else if len(rd) > 0 && (rd[0] == 'D' || rd == "H" || rd == "-H-") {
			p.RdScore = append(p.RdScore, 0.5)
			curTotal += 0.5
		} else {
			p.RdScore = append(p.RdScore, 0.0)
		}
		p.RdTotal = append(p.RdTotal, curTotal)
	}
}

// Parse populates a Player from a fixed-format rating report element slice
func (p *Player) Parse(elements []string, numRounds int) *Player {
	length := len(elements)
	idx := length - 1
	for i := 0; i < numRounds; i++ {
		rnd := fixRound(elements[idx-i])
		p.Rounds = append([]string{rnd}, p.Rounds...)
	}
	p.RatePost = elements[length-numRounds-1]
	p.RatePre = elements[length-numRounds-2]
	p.State = elements[length-numRounds-3]
	p.ID = elements[0]
	p.Rank = elements[1]
	p.Name = elements[2]
	p.State = elements[3]
	p.Name = utils.MassageName(p.Name)
	return p
}

func fixRound(r string) string {
	rnd := strings.ReplaceAll(r, "-", "")
	if len(rnd) == 2 && rnd[1] == '0' {
		return string(rnd[0])
	}
	return rnd
}

// ToHTML renders the player as a full crosstable table row (rating report format)
func (p *Player) ToHTML() string {
	var b strings.Builder
	b.WriteString("<tr>\n")
	b.WriteString("\t<td>" + p.Rank + "</td>\n")
	b.WriteString("\t<td>" + p.Name + "</td>\n")
	b.WriteString("\t<td>" + p.State + "</td>\n")
	b.WriteString("\t<td>" + p.ID + "</td>\n")
	b.WriteString("\t<td>" + p.RatePre + "</td>\n")
	b.WriteString("\t<td>" + p.RatePost + "</td>\n")
	b.WriteString("\t<td>" + p.Norm + "</td>\n")
	for _, r := range p.Rounds {
		b.WriteString("\t<td>" + r + "</td>\n")
	}
	b.WriteString("\t<td>" + p.Total + "</td>\n")
	b.WriteString("</tr>\n")
	return b.String()
}

// XtblHTML renders the player as a WinTD-style crosstable row
func (p *Player) XtblHTML(rowNbr, nbrRounds, scrRound int) string {
	var b strings.Builder
	b.WriteString("<tr>\n")
	b.WriteString(fmt.Sprintf("\t<td>%d</td>\n", rowNbr))
	b.WriteString("\t<td>" + p.Name + "</td>\n")
	b.WriteString("\t<td>" + p.RatePre + "</td>\n")
	for r := 0; r < nbrRounds; r++ {
		result := ""
		if r < len(p.Rounds) {
			result = p.Rounds[r]
		}
		b.WriteString("\t<td>" + result + "</td>\n")
	}
	total := 0.0
	if scrRound > 0 && scrRound <= len(p.RdTotal) {
		total = p.RdTotal[scrRound-1]
	}
	b.WriteString(fmt.Sprintf("\t<td>%g</td>\n", total))
	b.WriteString("</tr>\n")
	return b.String()
}

// GamesHTML renders a pairings row from parsed column elements
func GamesHTML(elements []string) string {
	get := func(i int) string {
		if i < len(elements) {
			return elements[i]
		}
		return ""
	}
	return "<tr>\n" +
		"\t<td>" + get(0) + "</td>\n" +
		"\t<td>" + get(1) + "</td>\n" +
		"\t<td>" + get(2) + "</td>\n" +
		"</tr>\n"
}

// ByeHTML renders a bye row from parsed column elements
func ByeHTML(elements []string) string {
	get := func(i int) string {
		if i < len(elements) {
			return elements[i]
		}
		return ""
	}
	return "<tr>\n" +
		"\t<td></td>\n" +
		"\t<td>" + get(1) + "</td>\n" +
		"\t<td>Please Wait</td>\n" +
		"</tr>\n"
}
