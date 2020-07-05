package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

const (
	lineBreak = "<br>"
)

func makeLink(source, content string) string {
	return fmt.Sprintf(`<a href="%v">%v</a>`, source, content)
}

func makeParagraph(content string) string {
	return `<p>` + content + `</p>`
}

func makeBlockQuote(content string) string {
	return `<blockquote>` + content + `</blockquote>`
}

func makeHeader(level int, content string) string {
	return fmt.Sprintf("<h%v>%v</h%v>", level, content, level)
}

func makePreformatted(content string) string {
	return "<pre>" + content + "</pre>"
}

func makeListItem(content string) string {
	return "<li>" + content + "</li>"
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
