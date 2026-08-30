package main

import (
	"bufio"
	"chess-utils/chess"
	"chess-utils/utils"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const EVENT_FLOOR_YEAR = 2019

var classCycle = []string{"wccColor1", "wccColor2", "wccColor3", "wccColor4", "wccColor5", "wccColor6"}

func getInputFilename() string {
	if len(os.Args) >= 3 {
		return os.Args[2]
	}
	return "../../resources/LateSpring-Open.txt"
}

// splitFullPostUpdateFile reads a WinTD combined report file and splits it into
// up to 4 sections (openXtbl, reserveXtbl, openPairings, reservePairings).
// A new section starts when a non-blank, non-indented line appears after content.
func splitFullPostUpdateFile() ([][]string, error) {
	filename := getInputFilename()
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var outLineSet [][]string
	var thisLineSet []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text() + "\n"
		if len(line) > 1 && line[0] != ' ' && len(thisLineSet) > 0 {
			outLineSet = append(outLineSet, thisLineSet)
			thisLineSet = []string{line}
		} else {
			thisLineSet = append(thisLineSet, line)
		}
	}
	if len(thisLineSet) > 0 {
		outLineSet = append(outLineSet, thisLineSet)
	}
	// Pad to 4 sections with nil
	for len(outLineSet) < 4 {
		outLineSet = append(outLineSet, nil)
	}
	return outLineSet, scanner.Err()
}

func fullUpdatePost() error {
	sections, err := splitFullPostUpdateFile()
	if err != nil {
		return err
	}
	openXtbl, reserveXtbl, openPair, reservePair := sections[0], sections[1], sections[2], sections[3]

	fmt.Println("</br><b>Open Section Crosstable</b></br>")
	processWinTDLines(os.Stdout, openXtbl)
	if openPair != nil {
		fmt.Println("</br><b>Open Section Pairings</b></br>")
		processGamesLines(os.Stdout, openPair)
	}
	fmt.Println("</br><b>Reserve Section Crosstable</b></br>")
	processWinTDLines(os.Stdout, reserveXtbl)
	if reservePair != nil {
		fmt.Println("</br><b>Reserve Section Pairings</b></br>")
		processGamesLines(os.Stdout, reservePair)
	}
	return nil
}

// processRatingReportFile handles the 'ratedEventFromRatingReport' command
func processRatingReportFile() error {
	filename := getInputFilename()
	numRounds := 4
	if len(os.Args) >= 4 {
		n, err := strconv.Atoi(os.Args[3])
		if err == nil {
			numRounds = n
		}
	}

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	utils.PrintTableHeader(os.Stdout, numRounds)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		elements := splitFixedLine(strings.TrimRight(line, "\r\n"), numRounds)
		p := chess.NewPlayer()
		p.Parse(elements, numRounds)
		fmt.Print(p.ToHTML())
	}
	fmt.Println("</tbody>")
	fmt.Println("</table>")
	utils.PrintPageClose(os.Stdout)
	return scanner.Err()
}

// processWinTDFile handles 'updatePostFromWinTdXTable' (reads from file)
func processWinTDFile() error {
	filename := getInputFilename()
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text()+"\n")
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	processWinTDLines(os.Stdout, lines)
	return nil
}

// processWinTDLines renders a WinTD crosstable section to w
func processWinTDLines(w *os.File, inLines []string) {
	numRounds := 4
	if len(os.Args) >= 4 {
		if n, err := strconv.Atoi(os.Args[3]); err == nil {
			numRounds = n
		}
	}
	utils.PrintCrossTableHeader(w, numRounds)

	players := 0
	for _, line := range inLines {
		if len(line) > 5 && line[4] == '.' && line[2] != 'N' {
			players++
			p := chess.CreatePlayerFromWinTDXtbl(line, numRounds)
			fmt.Fprint(w, p.XtblHTML(players, numRounds, numRounds))
		}
	}
	fmt.Fprintln(w, "</tbody>")
	fmt.Fprintln(w, "</table>")
}

// processGamesFile handles 'updatePostFromWinTdPairings' (reads from file)
func processGamesFile() error {
	filename := getInputFilename()
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text()+"\n")
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	processGamesLines(os.Stdout, lines)
	return nil
}

// processGamesLines renders WinTD pairings to w
func processGamesLines(w *os.File, inLines []string) {
	utils.PrintGamesTableHeader(w)
	for _, line := range inLines {
		elements := splitGamesLine(line)
		if len(line) > 7 && line[6] == '.' {
			fmt.Fprint(w, chess.GamesHTML(elements))
		} else if strings.Contains(line, "Please Wait") {
			fmt.Fprint(w, chess.ByeHTML(elements))
		}
	}
	fmt.Fprintln(w, "</tbody>")
	fmt.Fprintln(w, "</table>")
}

func processPastEvent() error {
	if len(os.Args) <= 2 {
		return fmt.Errorf("usage: chess-utils ratedEventByID <tournamentId>")
	}
	tournamentID := os.Args[2]
	r := chess.NewReader()
	tournament, err := r.GetPastEvent(tournamentID)
	if err != nil {
		return err
	}
	t := chess.NewTransformer()
	t.CreateBlogResults(os.Stdout, tournament)
	return nil
}

func addEventNameOverrides(events []*chess.Event, theClubAbbrev string) error {
	// Load name overrides
	filename := getEventOverridesFilename(theClubAbbrev)
	overrideData, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read %s: %w", filename, err)
	}
	var eventNameOverrides map[string]string
	if err := json.Unmarshal(overrideData, &eventNameOverrides); err != nil {
		return err
	}

	for _, event := range events {
		if overrideName, ok := eventNameOverrides[event.ID]; ok {
			event.SetNameOverride(overrideName)
		}
	}

	return nil
}

func getEventOverridesFilename(theClubAbbrev string) string {
	switch theClubAbbrev {
	case "wcc":
		return "data/wcc_event_names.json"
	case "swcc":
		return "data/swcc_event_names.json"
	default:
		return ""
	}
}

// processGenerateEventsJS handles the 'generateEventsJS' command: fetches all rated
// events for the given affiliate ID and writes data/web/events.js, a JS module
// exporting TOURNAMENTS grouped by year (newest year and newest event first).
func processGenerateEventsJS(theAffiliateId, theClubAbbrev string) error {
	if len(os.Args) <= 2 {
		return fmt.Errorf("usage: chess-utils generateEventsJS <clubAbbrev>")
	}

	fmt.Println("processing events for affiliate", theAffiliateId)
	r := chess.NewReader()
	events, err := r.GetPastEvents(theAffiliateId, EVENT_FLOOR_YEAR)
	if err != nil {
		return err
	}
	if len(events) == 0 {
		fmt.Println("No events returned for affiliate", theAffiliateId)
		return nil
	}

	err = addEventNameOverrides(events, theClubAbbrev)
	if err != nil {
		return err
	}

	foutName := "data/web/" + theClubAbbrev + "_events.js"
	fOut, err := os.Create(foutName)
	if err != nil {
		return err
	}
	defer fOut.Close()

	//   '2018 and earlier': [
	//    { name: 'View full archive on USCF →', url: 'https://ratings.uschess.org/affiliate/A6011047' },
	//  ],
	// Group events by year, preserving the newest-first order within each year
	var years []int
	byYear := make(map[int][]*chess.Event)

	for _, event := range events {
		year := event.FinishDate.Year()
		if _, ok := byYear[year]; !ok {
			years = append(years, year)
		}
		byYear[year] = append(byYear[year], event)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(years)))

	fmt.Fprintln(fOut, "export const TOURNAMENTS = {")
	for _, year := range years {
		fmt.Fprintf(fOut, "  '%d': [\n", year)
		for _, event := range byYear[year] {
			fmt.Fprintf(fOut, "    { name: %s, url: %s },\n", jsQuote(event.GetName()), jsQuote(event.Href))
		}
		fmt.Fprintln(fOut, "  ],")
	}

	url := "https://ratings.uschess.org/affiliate/" + theAffiliateId
	fmt.Fprintf(fOut, "  '%d': [\n", (EVENT_FLOOR_YEAR - 1))
	fmt.Fprintf(fOut, "    { name: %s, url: %s },\n", jsQuote("'View full archive on USCF →'"), jsQuote(url))
	fmt.Fprintln(fOut, "  ],")

	fmt.Fprintln(fOut, "};")

	// Also echo IDs + names to stdout for reference, same as clubEvents
	// commented out, because the output was being used to populate the 'preferred names' files, and that's not a good
	//for _, event := range events {
	//	outName := event.GetName()
	//	fmt.Printf("%q: %q,\n", event.ID, outName)
	//}
	return nil
}

// jsQuote wraps a string in single quotes, escaping backslashes and single
// quotes so it can be safely embedded as a JS string literal.
func jsQuote(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `'`, `\'`)
	return "'" + s + "'"
}

// splitFixedLine parses a fixed-width rating report line
//
// 012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789
// 30289677    1 KLINKNER, PATRICK WI  927/24  921*   W---7 L---5 W---3 W---2 W---2 3.0
func splitFixedLine(line string, numRounds int) []string {
	if line == "" {
		return nil
	}
	// Pad if needed
	for len(line) < 51+(6*numRounds) {
		line += " "
	}
	elements := []string{
		strings.TrimSpace(line[0:9]),
		strings.TrimSpace(line[9:13]),
		strings.TrimSpace(line[14:31]),
		strings.TrimSpace(line[32:34]),
		strings.TrimSpace(line[35:43]),
		strings.TrimSpace(line[43:51]),
	}
	for i := 0; i < numRounds; i++ {
		start := 51 + (6 * i)
		end := start + 6
		if end > len(line) {
			end = len(line)
		}
		elements = append(elements, strings.TrimSpace(line[start:end]))
	}
	return elements
}

// splitGamesLine parses a fixed-width WinTD pairings line
//
// 012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789
//
//  101. ___  Templin, Aethe (2.0,Templi,1960)  ___  Coons, James Jay (2.0,1724)
func splitGamesLine(line string) []string {
	if line == "" {
		return nil
	}
	for len(line) < 90 {
		line += " "
	}
	return []string{
		strings.TrimSpace(line[3:7]),
		strings.TrimSpace(line[14:47]),
		strings.TrimSpace(line[53:90]),
	}
}

func usage() {
	fmt.Println(`Usage: chess-utils <option> [param]

Options:
  postUpdate <file>                  Weekly update post from WinTD combined report
  ratedEventByID <eventID>           Blog HTML for a completed rated event
  ratedEventFromRatingReport <file>  Blog HTML from USCF rating report text file
  updatePostFromWinTdXTable <file>   Crosstable HTML from one WinTD section file
  updatePostFromWinTdPairings <file> Pairings HTML from one WinTD section file
  generateEventsJS <clubAbbrev>      Generate data/web/events.js; arg is either wcc or swcc
  clubEvents                         [DEPRECATED] Use generateEventsJS instead
  winnersPage                        [DEPRECATED] There is no replacement, since past champions are in a JavaScript array
  ratedEventFromFile <file>          [DEPRECATED] Use ratedEventByID instead`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "postUpdate":
		err = fullUpdatePost()
	case "ratedEventByID":
		err = processPastEvent()
	case "ratedEventFromFile":
		fmt.Println("The option 'ratedEventFromFile' should not be used anymore.")
		fmt.Println("Please use 'ratedEventByID'.")
		os.Exit(1)
	case "ratedEventFromRatingReport":
		err = processRatingReportFile()
	case "updatePostFromWinTdXTable":
		err = processWinTDFile()
	case "updatePostFromWinTdPairings":
		err = processGamesFile()
	case "clubEvents":
		fmt.Println("The option 'clubEvents' should not be used anymore.")
		fmt.Println("Please use 'generateEventsJS'.")
		os.Exit(1)
	case "winnersPage":
		fmt.Println("The option 'winnersPage' should not be used anymore.")
		fmt.Println("There is no direct replacement'.")
		os.Exit(1)
	case "generateEventsJS":
		if len(os.Args) >= 3 {
			var clubAbbrev string
			var affiliateId string
			switch os.Args[2] {
			case "wcc":
				clubAbbrev = os.Args[2]
				affiliateId = chess.WccAffiliateID
			case "swcc":
				clubAbbrev = os.Args[2]
				affiliateId = chess.SwccAffiliateID
			}
			err = processGenerateEventsJS(affiliateId, clubAbbrev)
		} else {
			usage()
			os.Exit(1)
		}
	default:
		usage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
