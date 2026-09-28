package utils

import (
	"encoding/csv"
	"fmt"
	"html"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

// ParseSeparator maps a separator argument to a rune. Accepts named aliases
// (comma, tab, pipe, semicolon) or any literal single character. Empty means comma.
func ParseSeparator(arg string) (rune, error) {
	switch strings.ToLower(arg) {
	case "", "comma":
		return ',', nil
	case "tab":
		return '\t', nil
	case "pipe":
		return '|', nil
	case "semicolon":
		return ';', nil
	}
	if utf8.RuneCountInString(arg) == 1 {
		r, _ := utf8.DecodeRuneInString(arg)
		if r == '"' || r == '\r' || r == '\n' {
			return 0, fmt.Errorf("invalid separator %q", arg)
		}
		return r, nil
	}
	return 0, fmt.Errorf("unknown separator %q (use comma, tab, pipe, semicolon, or a single character)", arg)
}

// ConvertCsvToHTML reads delimited records from r and writes a wccCrosstable
// HTML table fragment to w. The first record becomes the header row. Quoted
// values (including embedded separators, "" escapes and newlines) are handled
// by encoding/csv; surrounding quotes are dropped. Values are HTML-escaped,
// trimmed, and short rows are padded to the widest row.
func ConvertCsvToHTML(r io.Reader, w io.Writer, sep rune) error {
	cr := csv.NewReader(r)
	cr.Comma = sep
	cr.FieldsPerRecord = -1 // allow ragged rows
	cr.LazyQuotes = true    // tolerate stray quotes in unquoted fields
	// Only skip leading space when the separator itself isn't whitespace;
	// otherwise leading tabs (empty leading columns) would be swallowed.
	// Cell values are still trimmed in printCsvRow.
	cr.TrimLeadingSpace = !unicode.IsSpace(sep)

	records, err := cr.ReadAll() // blank lines are skipped by encoding/csv
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return nil
	}

	width := 0
	for _, rec := range records {
		if len(rec) > width {
			width = len(rec)
		}
	}

	fmt.Fprintln(w, "<table class='wccCrosstable'>")
	fmt.Fprintln(w, "<thead>")
	printCsvRow(w, records[0], width, "th")
	fmt.Fprintln(w, "</thead>")
	fmt.Fprintln(w, "<tbody>")
	for _, rec := range records[1:] {
		printCsvRow(w, rec, width, "td")
	}
	fmt.Fprintln(w, "</tbody>")
	fmt.Fprintln(w, "</table>")
	return nil
}

func printCsvRow(w io.Writer, rec []string, width int, tag string) {
	fmt.Fprintln(w, "<tr>")
	for i := 0; i < width; i++ {
		val := ""
		if i < len(rec) {
			val = html.EscapeString(strings.TrimSpace(rec[i]))
		}
		fmt.Fprintf(w, "\t<%s>%s</%s>\n", tag, val, tag)
	}
	fmt.Fprintln(w, "</tr>")
}
