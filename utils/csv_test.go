package utils

import (
	"strings"
	"testing"
)

func TestParseSeparator(t *testing.T) {
	cases := map[string]rune{"": ',', "comma": ',', "TAB": '\t', "pipe": '|', "semicolon": ';', ";": ';', "|": '|'}
	for in, want := range cases {
		got, err := ParseSeparator(in)
		if err != nil || got != want {
			t.Errorf("ParseSeparator(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
	for _, bad := range []string{`"`, "bogus"} {
		if _, err := ParseSeparator(bad); err == nil {
			t.Errorf("ParseSeparator(%q) expected error", bad)
		}
	}
}

func TestConvertCsvToHTML(t *testing.T) {
	in := "Name,Rating,Club\n" +
		"\"Smith, John\",1850,WCC\n" +
		"\n" +
		"  Doe ,1700\n" +
		"\"Say \"\"Hi\"\"\",<b>,A&B\n"
	var sb strings.Builder
	if err := ConvertCsvToHTML(strings.NewReader(in), &sb, ','); err != nil {
		t.Fatal(err)
	}
	out := sb.String()
	for _, want := range []string{
		"<table class='wccCrosstable'>",
		"\t<th>Name</th>",
		"\t<td>Smith, John</td>",
		"\t<td>Doe</td>\n\t<td>1700</td>\n\t<td></td>", // trimmed + padded
		"\t<td>Say &#34;Hi&#34;</td>",
		"\t<td>&lt;b&gt;</td>",
		"\t<td>A&amp;B</td>",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in output:\n%s", want, out)
		}
	}
	if strings.Count(out, "<tr>") != 4 {
		t.Errorf("expected 4 rows (blank line skipped), got:\n%s", out)
	}
}

func TestConvertCsvToHTMLPipe(t *testing.T) {
	var sb strings.Builder
	if err := ConvertCsvToHTML(strings.NewReader("a|b\n\"x|y\"|z\n"), &sb, '|'); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(sb.String(), "<td>x|y</td>") {
		t.Errorf("quoted pipe not handled:\n%s", sb.String())
	}
}

func TestConvertCsvToHTMLLeadingTabs(t *testing.T) {
	in := "A\tB\tC\n\tx\ty\n\t\tz\n"
	var sb strings.Builder
	if err := ConvertCsvToHTML(strings.NewReader(in), &sb, '\t'); err != nil {
		t.Fatal(err)
	}
	out := sb.String()
	for _, want := range []string{
		"<tr>\n\t<td></td>\n\t<td>x</td>\n\t<td>y</td>\n</tr>",
		"<tr>\n\t<td></td>\n\t<td></td>\n\t<td>z</td>\n</tr>",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("leading separators lost; missing %q in:\n%s", want, out)
		}
	}
}

func TestConvertCsvToHTMLLeadingCommas(t *testing.T) {
	var sb strings.Builder
	if err := ConvertCsvToHTML(strings.NewReader("A,B,C\n,,z\n"), &sb, ','); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(sb.String(), "<tr>\n\t<td></td>\n\t<td></td>\n\t<td>z</td>\n</tr>") {
		t.Errorf("leading commas lost:\n%s", sb.String())
	}
}
