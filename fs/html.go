package fs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/russross/blackfriday.v2"
)

func markdownToHtml(md []byte) string {
	return string(blackfriday.Run(NormalizeEOL(md)))
}

func (h *Hypha) asHtml() (string, error) {
	rev := h.actual
	ret := `<article class="page">
	<h1 class="page__title">` + rev.FullName + `</h1>
`
	// What about using <figure>?
	if h.hasBinaryData() {
		ret += fmt.Sprintf(`<img src="/%s?action=binary&rev=%d" class="page__amnt"/>`, rev.FullName, rev.Id)
	}

	contents, err := ioutil.ReadFile(rev.TextPath)
	if err != nil {
		log.Println("Failed to render", rev.FullName, ":", err)
		return "", err
	}

	// TODO: support more markups.
	// TODO: support mycorrhiza extensions like transclusion.
	switch rev.TextMime {
	case "text/markdown":
		html := markdownToHtml(contents)
		ret += string(html)
	default:
		ret += fmt.Sprintf(`<pre>%s</pre>`, contents)
	}

	ret += `
</article>`

	return ret, nil
}

// NormalizeEOL will convert Windows (CRLF) and Mac (CR) EOLs to UNIX (LF)
// Code taken from here: https://github.com/go-gitea/gitea/blob/dc8036dcc680abab52b342d18181a5ee42f40318/modules/util/util.go#L68-L102
// Gitea has MIT License
//
// We use it because md parser does not handle CRLF correctly. I don't know why, but CRLF appears sometimes.
func NormalizeEOL(input []byte) []byte {
	var right, left, pos int
	if right = bytes.IndexByte(input, '\r'); right == -1 {
		return input
	}
	length := len(input)
	tmp := make([]byte, length)

	// We know that left < length because otherwise right would be -1 from IndexByte.
	copy(tmp[pos:pos+right], input[left:left+right])
	pos += right
	tmp[pos] = '\n'
	left += right + 1
	pos++

	for left < length {
		if input[left] == '\n' {
			left++
		}

		right = bytes.IndexByte(input[left:], '\r')
		if right == -1 {
			copy(tmp[pos:], input[left:])
			pos += length - left
			break
		}
		copy(tmp[pos:pos+right], input[left:left+right])
		pos += right
		tmp[pos] = '\n'
		left += right + 1
		pos++
	}
	return tmp[:pos]
}
