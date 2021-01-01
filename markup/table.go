package markup

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	// "github.com/bouncepaw/mycorrhiza/util"
)

var tableRe = regexp.MustCompile(`^table\s+{`)

func MatchesTable(line string) bool {
	return tableRe.MatchString(line)
}

func TableFromFirstLine(line, hyphaName string) *Table {
	return &Table{
		hyphaName: hyphaName,
		caption:   line[strings.IndexRune(line, '{')+1:],
		rows:      make([]*tableRow, 0),
	}
}

func (t *Table) Process(line string) (shouldGoBackToNormal bool) {
	if strings.TrimSpace(line) == "}" && !t.inMultiline {
		return true
	}
	if !t.inMultiline {
		t.pushRow()
	}
	var (
		escaping           bool
		lookingForNonSpace = !t.inMultiline
		countingColspan    bool
	)
	for i, r := range line {
		switch {
		case lookingForNonSpace && unicode.IsSpace(r):
		case lookingForNonSpace && (r == '!' || r == '|'):
			t.currCellMarker = r
			t.currColspan = 1
			lookingForNonSpace = false
			countingColspan = true
		case lookingForNonSpace:
			t.currCellMarker = '^' // ^ represents implicit |, not part of syntax
			t.currColspan = 1
			lookingForNonSpace = false
			t.currCellBuilder.WriteRune(r)

		case escaping:
			t.currCellBuilder.WriteRune(r)

		case t.inMultiline && r == '}':
			t.inMultiline = false
		case t.inMultiline && i == len(line)-1:
			t.currCellBuilder.WriteRune('\n')
		case t.inMultiline:
			t.currCellBuilder.WriteRune(r)

			// Not in multiline:
		case (r == '|' || r == '!') && !countingColspan:
			t.pushCell()
			t.currCellMarker = r
			t.currColspan = 1
			countingColspan = true
		case r == t.currCellMarker && (r == '|' || r == '!') && countingColspan:
			t.currColspan++
		case r == '{':
			t.inMultiline = true
			countingColspan = false
		case i == len(line)-1:
			t.pushCell()
		default:
			t.currCellBuilder.WriteRune(r)
			countingColspan = false
		}
	}
	return false
}

type Table struct {
	// data
	hyphaName string
	caption   string
	rows      []*tableRow
	// state
	inMultiline bool
	// tmp
	currCellMarker  rune
	currColspan     uint
	currCellBuilder strings.Builder
}

func (t *Table) pushRow() {
	t.rows = append(t.rows, &tableRow{
		cells: make([]*tableCell, 0),
	})
}

func (t *Table) pushCell() {
	tc := &tableCell{
		content: t.currCellBuilder.String(),
		colspan: t.currColspan,
	}
	switch t.currCellMarker {
	case '|', '^':
		tc.kind = tableCellDatum
	case '!':
		tc.kind = tableCellHeader
	}
	// We expect the table to have at least one row ready, so no nil-checking
	tr := t.rows[len(t.rows)-1]
	tr.cells = append(tr.cells, tc)
	t.currCellBuilder = strings.Builder{}
}

func (t *Table) asHtml() (html string) {
	if t.caption != "" {
		html += fmt.Sprintf("<caption>%s</caption>", t.caption)
	}
	if len(t.rows) > 0 && t.rows[0].looksLikeThead() {
		html += fmt.Sprintf("<thead>%s</thead>", t.rows[0].asHtml(t.hyphaName))
		t.rows = t.rows[1:]
	}
	html += "\n<tbody>\n"
	for _, tr := range t.rows {
		html += tr.asHtml(t.hyphaName)
	}
	return fmt.Sprintf(`<table>%s</tbody></table>`, html)
}

type tableRow struct {
	cells []*tableCell
}

func (tr *tableRow) asHtml(hyphaName string) (html string) {
	for _, tc := range tr.cells {
		html += tc.asHtml(hyphaName)
	}
	return fmt.Sprintf("<tr>%s</tr>\n", html)
}

// Most likely, rows with more than two header cells are theads. I allow one extra datum cell for tables like this:
// |   ! a ! b
// ! c | d | e
// ! f | g | h
func (tr *tableRow) looksLikeThead() bool {
	var (
		headerAmount = 0
		datumAmount  = 0
	)
	for _, tc := range tr.cells {
		switch tc.kind {
		case tableCellHeader:
			headerAmount++
		case tableCellDatum:
			datumAmount++
		}
	}
	return headerAmount >= 2 && datumAmount <= 1
}

type tableCell struct {
	kind    tableCellKind
	colspan uint
	content string
}

func (tc *tableCell) asHtml(hyphaName string) string {
	return fmt.Sprintf(
		"<%[1]s %[2]s>%[3]s</%[1]s>\n",
		tc.kind.tagName(),
		tc.colspanAttribute(),
		tc.contentAsHtml(hyphaName),
	)
}

func (tc *tableCell) colspanAttribute() string {
	if tc.colspan <= 1 {
		return ""
	}
	return fmt.Sprintf(`colspan="%d"`, tc.colspan)
}

func (tc *tableCell) contentAsHtml(hyphaName string) (html string) {
	for _, line := range strings.Split(tc.content, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			if html != "" {
				html += `<br>`
			}
			html += ParagraphToHtml(hyphaName, line)
		}
	}
	return html
}

type tableCellKind int

const (
	tableCellUnknown tableCellKind = iota
	tableCellHeader
	tableCellDatum
)

func (tck tableCellKind) tagName() string {
	switch tck {
	case tableCellHeader:
		return "th"
	case tableCellDatum:
		return "td"
	default:
		return "p"
	}
}
