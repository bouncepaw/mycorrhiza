package markup

import (
	"bytes"
	"fmt"
	"html"
	"regexp"
	"strings"
	"unicode"
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
	spanStrike
	spanLink
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

func getLinkNode(input *bytes.Buffer, hyphaName string, isBracketedLink bool) string {
	if isBracketedLink {
		input.Next(2) // drop those [[
	}
	var (
		escaping   = false
		addrBuf    = bytes.Buffer{}
		displayBuf = bytes.Buffer{}
		currBuf    = &addrBuf
	)
	for input.Len() != 0 {
		b, _ := input.ReadByte()
		if escaping {
			currBuf.WriteByte(b)
			escaping = false
		} else if isBracketedLink && b == '|' && currBuf == &addrBuf {
			currBuf = &displayBuf
		} else if isBracketedLink && b == ']' && bytes.HasPrefix(input.Bytes(), []byte{']'}) {
			input.Next(1)
			break
		} else if !isBracketedLink && unicode.IsSpace(rune(b)) {
			break
		} else {
			currBuf.WriteByte(b)
		}
	}
	href, text, class := LinkParts(addrBuf.String(), displayBuf.String(), hyphaName)
	return fmt.Sprintf(`<a href="%s" class="%s">%s</a>`, href, class, html.EscapeString(text))
}

// getTextNode splits the `input` into two parts `textNode` and `rest` by the first encountered rune that resembles a span tag. If there is none, `textNode = input`, `rest = ""`. It handles escaping with backslash.
func getTextNode(input *bytes.Buffer) string {
	var (
		textNodeBuffer = bytes.Buffer{}
		escaping       = false
		startsWith     = func(t string) bool {
			return bytes.HasPrefix(input.Bytes(), []byte(t))
		}
		couldBeLinkStart = func() bool {
			return startsWith("https://") || startsWith("http://") || startsWith("gemini://") || startsWith("gopher://") || startsWith("ftp://")
		}
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
		} else if strings.IndexByte("/*`^,![~", b) >= 0 {
			input.UnreadByte()
			break
		} else if couldBeLinkStart() {
			textNodeBuffer.WriteByte(b)
			break
		} else {
			textNodeBuffer.WriteByte(b)
		}
	}
	return textNodeBuffer.String()
}

var (
	dangerousSymbols = "<>{}|\\^[]`,()"
	reLink           = regexp.MustCompile(fmt.Sprintf(`[^[]{0,2}((https|http|gemini|gopher)://[^%[1]s]+)|(mailto:[^%[1]s]+)[^]]{0,2}`, dangerousSymbols))
)

// TODO:
func doRegexpStuff(input string) string {
	reLink.ReplaceAllString(input, "[[$1]]")
	return ""
}

func ParagraphToHtml(hyphaName, input string) string {
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
			spanLink:   false,
		}
		startsWith = func(t string) bool {
			return bytes.HasPrefix(p.Bytes(), []byte(t))
		}
		noTagsActive = func() bool {
			return !(tagState[spanItalic] || tagState[spanBold] || tagState[spanMono] || tagState[spanSuper] || tagState[spanSub] || tagState[spanMark] || tagState[spanLink])
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
		case startsWith("~~"):
			ret.WriteString(tagFromState(spanMark, tagState, "s", "~~"))
			p.Next(2)
		case startsWith("[["):
			ret.WriteString(getLinkNode(p, hyphaName, true))
		case (startsWith("https://") || startsWith("http://") || startsWith("gemini://") || startsWith("gopher://") || startsWith("ftp://")) && noTagsActive():
			ret.WriteString(getLinkNode(p, hyphaName, false))
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
			case spanStrike:
				ret.WriteString(tagFromState(spanMark, tagState, "s", "~~"))
			case spanLink:
				ret.WriteString(tagFromState(spanLink, tagState, "a", "[["))
			}
		}
	}

	return ret.String()
}
