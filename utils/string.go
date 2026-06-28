package utils

import (
	"fmt"
	"io"
	"strings"
)

// MassageName converts a name to title case with special McC handling
func MassageName(name string) string {
	result := ""
	cnc := true
	for _, c := range strings.ToLower(name) {
		if cnc {
			result += strings.ToUpper(string(c))
			cnc = false
		} else {
			result += string(c)
		}
		if c == ' ' || c == '\'' {
			cnc = true
		}
	}
	if strings.HasPrefix(result, "Mcc") {
		result = strings.Replace(result, "Mcc", "McC", 1)
	}
	return result
}

func PrintTableHeader(w io.Writer, numRounds int) {
	fmt.Fprintln(w, "<table class='wccCrosstable'>")
	fmt.Fprintln(w, "<thead>")
	fmt.Fprintln(w, "<tr>")
	fmt.Fprintln(w, "\t<th>No.</th>")
	fmt.Fprintln(w, "\t<th>Player Name</th>")
	fmt.Fprintln(w, "\t<th>State</th>")
	fmt.Fprintln(w, "\t<th>USCF ID</th>")
	fmt.Fprintln(w, "\t<th>Pre</th>")
	fmt.Fprintln(w, "\t<th>Post</th>")
	fmt.Fprintln(w, "\t<th>Norm</th>")
	for i := 0; i < numRounds; i++ {
		fmt.Fprintf(w, "\t<th>R%d</th>\n", i+1)
	}
	fmt.Fprintln(w, "\t<th>Total</th>")
	fmt.Fprintln(w, "</tr>")
	fmt.Fprintln(w, "</thead>")
	fmt.Fprintln(w, "<tbody>")
}

// PrintCrossTableHeader prints the WinTD-style crosstable header (no norm column)
func PrintCrossTableHeader(w io.Writer, numRounds int) {
	fmt.Fprintln(w, "<table class='wccCrosstable'>")
	fmt.Fprintln(w, "<thead>")
	fmt.Fprintln(w, "<tr>")
	fmt.Fprintln(w, "\t<th>No.</th>")
	fmt.Fprintln(w, "\t<th>Player Name</th>")
	fmt.Fprintln(w, "\t<th>Rating</th>")
	for i := 0; i < numRounds; i++ {
		fmt.Fprintf(w, "\t<th>R%d</th>\n", i+1)
	}
	fmt.Fprintln(w, "\t<th>Total</th>")
	fmt.Fprintln(w, "</tr>")
	fmt.Fprintln(w, "</thead>")
	fmt.Fprintln(w, "<tbody>")
}

func PrintGamesTableHeader(w io.Writer) {
	fmt.Fprintln(w, "<table class='wccCrosstable'>")
	fmt.Fprintln(w, "<thead>")
	fmt.Fprintln(w, "<tr>")
	fmt.Fprintln(w, "\t<th>No.</th>")
	fmt.Fprintln(w, "\t<th>White</th>")
	fmt.Fprintln(w, "\t<th>Black</th>")
	fmt.Fprintln(w, "</tr>")
	fmt.Fprintln(w, "</thead>")
	fmt.Fprintln(w, "<tbody>")
}

func PrintPageHeader(w io.Writer, title, subtitle string) {
	fmt.Fprintf(w, `<! DOCTYPE html>
    <html><head>
    <meta http-equiv=Content-Type content="text/html" charset="utf-8">
    <link rel="stylesheet" type="text/css" href="css/wcc.css"/>
    <script type="text/javascript" language="javascript" src="js/jquery-3.4.1.min.js">
        <!-- yes, this comment is here on purpose -->
    </script>
    <title>Waukesha Chess Club</title>
    </head>

    <body>
    <div id="main-container">
    <script>
        const isDarkMode = window.matchMedia("(prefers-color-scheme: dark)").matches;
        $(function() {
            if (isDarkMode)
            {
                $("#title").load("headerw.html");
            } else {
              $("#title").load("headerb.html");
            }
            $("#menu").load("menu.html");
            $("#footer").load("footer.html");
        });
    </script>
    <div id="title"> </div>
    <div id="menu"></div>
    <div id="divContent">
    <div id="pageTitle">
    <span>%s<br/></span>
    <span id="divSubtitle">%s</span></div>
`, title, subtitle)
}

func PrintEventTableHeader(w io.Writer, headers []string) {
	fmt.Fprintf(w, "<table class='wccPast' width=\"90%%\">\n<tr>")
	pct := 80 / len(headers)
	for _, hdr := range headers {
		fmt.Fprintf(w, "<th style=\"width: %d%%\" class =\"thPast\">%s</th>\n", pct, hdr)
	}
	fmt.Fprintln(w, "</tr>")
}

func PrintPriorEventRow(w io.Writer, newYear int, name, link, className string) {
	const aggregateFromYear = 1991
	if newYear > 0 && newYear >= aggregateFromYear {
		fmt.Fprintln(w, "<tr>\t<td></td> <td></td> </tr>")
		fmt.Fprintln(w, "<tr><td colspan=\"2\"><hr width=\"100%\" /></td></tr>")
		fmt.Fprintln(w, "<tr>")
		if newYear > aggregateFromYear {
			fmt.Fprintf(w, "<td class=\"year\">%d</td>\n", newYear)
		} else {
			fmt.Fprintf(w, "<td class=\"year\">%d (and earlier)</td>\n", newYear)
		}
	} else {
		fmt.Fprintln(w, "<tr>")
		fmt.Fprintln(w, "<td></td>")
	}
	fmt.Fprintf(w, "<td class='%s'><a class='%s' target=\"_blank\" href=\"%s\" >%s</a><br/></td>\n",
		className, className, link, name)
	fmt.Fprintln(w, "</tr>")
}

func PrintWinnersRow(w io.Writer, newYear int, name1, link1, name2, link2, className string) {
	if newYear > 0 {
		fmt.Fprintln(w, "<tr>\t<td></td> <td></td> </tr>")
		fmt.Fprintln(w, "<tr><td colspan=\"3\"><hr width=\"100%\" /></td></tr>")
		fmt.Fprintln(w, "<tr>")
		fmt.Fprintf(w, "<td class=\"year\">%d</td>\n", newYear)
	} else {
		fmt.Fprintln(w, "<tr>")
		fmt.Fprintln(w, "<td></td>")
	}
	if link1 == "" {
		fmt.Fprintf(w, "<td class='%s'>%s<br/></td>\n", className, name1)
	} else {
		fmt.Fprintf(w, "<td class='%s'><a class='%s' target=\"_blank\" href=\"%s\">%s</a><br/></td>\n",
			className, className, link1, name1)
	}
	if link2 == "" {
		fmt.Fprintf(w, "<td class='%s'>%s<br/></td>\n", className, name2)
	} else {
		fmt.Fprintf(w, "<td class='%s'><a class='%s' target=\"_blank\" href=\"%s\">%s</a><br/></td>\n",
			className, className, link2, name2)
	}
	fmt.Fprintln(w, "</tr>")
}

func PrintBlankLine(w io.Writer) {
	fmt.Fprintln(w, "<tr><td></td></tr>")
}

func PrintPageClose(w io.Writer) {
	fmt.Fprintln(w, "</div>")
	fmt.Fprintln(w, "<div id=\"footer\"></div>")
	fmt.Fprintln(w, "</body>")
	fmt.Fprintln(w, "</html>")
}

func PrintDivClose(w io.Writer) {
	fmt.Fprintln(w, "</div>")
}

func PrintTableClose(w io.Writer) {
	fmt.Fprintln(w, "</table>")
	fmt.Fprintln(w, "<br/>")
}
