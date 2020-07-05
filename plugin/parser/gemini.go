package parser

import (
	"bytes"
	"strings"

	"github.com/bouncepaw/mycorrhiza/util"
)

const (
	linkToken         = "=>"
	headerToken       = "#"
	quoteToken        = ">"
	preformattedToken = "```"
	listItemToken     = "*"
)

func GeminiToHtml(gemini []byte) string {
	lines, _ := StringToLines(string(util.NormalizeEOL(gemini)))
	var html []string
	for _, line := range lines {
		html = append(html, geminiLineToHtml(line))
	}
	buffer := bytes.Buffer{}
	for _, line := range html {
		buffer.WriteString(line)
	}
	return buffer.String()
}

func geminiLineToHtml(line string) (res string) {
	arr := strings.Fields(line)
	if len(arr) == 0 {
		return lineBreak
	}

	content := arr[1:]
	token := arr[0]
	if string(token[0]) == headerToken {
		return makeHeader(makeOutHeader(arr))
	}

	switch token {
	case linkToken:
		res = makeLink(makeOutLink(content))
	case quoteToken:
		res = makeBlockQuote(LinesToString(content, " "))
	case preformattedToken:
		res = makePreformatted(LinesToString(content, " "))
	case listItemToken:
		res = makeListItem(LinesToString(content, " "))
	default:
		res = makeParagraph(line)
	}
	return res
}

func makeOutLink(arr []string) (source, content string) {
	switch len(arr) {
	case 0:
		return "", ""
	case 1:
		return arr[0], arr[0]
	default:
		return arr[0], LinesToString(arr[1:], " ")
	}
}

func makeOutHeader(arr []string) (level int, content string) {
	level = len(arr[0])
	content = LinesToString(arr[1:], " ")
	return
}
