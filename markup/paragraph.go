package markup

import (
	"bytes"
	"fmt"
	"html"
	"strings"
)

type spanTokenType int

const (
	spanTextNode = iota
	spanItalic
	spanBold
	spanMono
	spanSuper
	spanSub
	spanMark
)

func tagFromState(stt spanTokenType, tagState map[spanTokenType]bool, tagName, originalForm string) string {
	if tagState[spanMono] && (stt != spanMono) {
		return originalForm
	}
	if tagState[stt] {
		tagState[stt] = false
		return fmt.Sprintf("</%s>", tagName)
	} else {
		tagState[stt] = true
		return fmt.Sprintf("<%s>", tagName)
	}
}

// getTextNode splits the `p` into two parts `textNode` and `rest` by the first encountered rune that resembles a span tag. If there is none, `textNode = p`, `rest = ""`. It handles escaping with backslash.
func getTextNode(input *bytes.Buffer) string {
	var (
		textNodeBuffer = bytes.Buffer{}
		escaping       = false
	)
	// Always read the first byte in advance to avoid endless loops that kill computers (sad experience)
	if input.Len() != 0 {
		b, _ := input.ReadByte()
		textNodeBuffer.WriteByte(b)
	}
	for input.Len() != 0 {
		// Assume no error is possible because we check for length
		b, _ := input.ReadByte()
		if escaping {
			textNodeBuffer.WriteByte(b)
			escaping = false
		} else if b == '\\' {
			escaping = true
		} else if strings.IndexByte("/*`^,!", b) >= 0 {
			input.UnreadByte()
			break
		} else {
			textNodeBuffer.WriteByte(b)
		}
	}
	return textNodeBuffer.String()
}

func ParagraphToHtml(input string) string {
	var (
		p   = bytes.NewBufferString(input)
		ret strings.Builder
		// true = tag is opened, false = tag is not opened
		tagState = map[spanTokenType]bool{
			spanItalic: false,
			spanBold:   false,
			spanMono:   false,
			spanSuper:  false,
			spanSub:    false,
			spanMark:   false,
		}
		startsWith = func(t string) bool {
			return bytes.HasPrefix(p.Bytes(), []byte(t))
		}
	)

	for p.Len() != 0 {
		switch {
		case startsWith("//"):
			ret.WriteString(tagFromState(spanItalic, tagState, "em", "//"))
			p.Next(2)
		case startsWith("**"):
			ret.WriteString(tagFromState(spanBold, tagState, "strong", "**"))
			p.Next(2)
		case startsWith("`"):
			ret.WriteString(tagFromState(spanMono, tagState, "code", "`"))
			p.Next(1)
		case startsWith("^"):
			ret.WriteString(tagFromState(spanSuper, tagState, "sup", "^"))
			p.Next(1)
		case startsWith(",,"):
			ret.WriteString(tagFromState(spanSub, tagState, "sub", ",,"))
			p.Next(2)
		case startsWith("!!"):
			ret.WriteString(tagFromState(spanMark, tagState, "mark", "!!"))
			p.Next(2)
		default:
			ret.WriteString(html.EscapeString(getTextNode(p)))
		}
	}

	for stt, open := range tagState {
		if open {
			switch stt {
			case spanItalic:
				ret.WriteString(tagFromState(spanItalic, tagState, "em", "//"))
			case spanBold:
				ret.WriteString(tagFromState(spanBold, tagState, "strong", "**"))
			case spanMono:
				ret.WriteString(tagFromState(spanMono, tagState, "code", "`"))
			case spanSuper:
				ret.WriteString(tagFromState(spanSuper, tagState, "sup", "^"))
			case spanSub:
				ret.WriteString(tagFromState(spanSub, tagState, "sub", ",,"))
			case spanMark:
				ret.WriteString(tagFromState(spanMark, tagState, "mark", "!!"))
			}
		}
	}

	return ret.String()
}
