package parser

import (
	"bufio"
	"bytes"
	"fmt"
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

var preState bool
var listState bool

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
	token := checkLineType(arr)

	switch token {
	case headerToken:
		level, content := makeOutHeader(arr)
		res = fmt.Sprintf("<h%v>%v</h%v>", level, content, level)
	case linkToken:
		source, content := makeOutLink(arr[1:])
		res = fmt.Sprintf(`<p><a href="%v">%v</a></p>`, source, content)
	case quoteToken:
		res = "<blockquote>" + LinesToString(arr[1:], " ") + "</blockquote>"
	case preformattedToken:
		preState = true
		res = fmt.Sprintf(`<pre alt="%v">`, LinesToString(arr[1:], " "))
	case "pre/empty":
		res = "\n"
	case "pre/text":
		res = line + "\n"
	case "pre/end":
		preState = false
		res = "</pre>"
	case "list/begin":
		res = "<ul><li>" + LinesToString(arr[1:], " ") + "</li>"
	case listItemToken:
		res = "<li>" + LinesToString(arr[1:], " ") + "</li>"
	case "list/end":
		listState = false
		res = "</ul>" + geminiLineToHtml(line)
	case "linebreak":
		res = "<br>"
	default:
		res = "<p>" + line + "</p>"
	}
	return
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

func checkLineType(arr []string) (res string) {
	isEmpty := len(arr) == 0
	if preState {
		if isEmpty {
			res = "pre/empty"
		} else if arr[0] == preformattedToken {
			res = "pre/end"
		} else {
			res = "pre/text"
		}
	} else if listState {
		if arr[0] == listItemToken {
			res = listItemToken
		} else {
			res = "list/end"
		}
	} else if isEmpty {
		res = "linebreak"
	} else if arr[0][0] == headerToken[0] {
		res = headerToken
	} else {
		return arr[0]
	}
	return
}

func StringToLines(s string) (lines []string, err error) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	return
}

func LinesToString(lines []string, separator string) string {
	buffer := bytes.Buffer{}
	for _, line := range lines {
		buffer.WriteString(line + separator)
	}
	return buffer.String()
}
