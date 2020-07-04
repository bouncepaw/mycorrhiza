package util

import (
	"bytes"
	"strings"
	"unicode"
)

func addColonPerhaps(name string) string {
	if strings.HasPrefix(name, ":") {
		return name
	}
	return ":" + name
}

func removeColonPerhaps(name string) string {
	if strings.HasPrefix(name, ":") {
		return name[1:]
	}
	return name
}

func UrlToCanonical(name string) string {
	return removeColonPerhaps(
		strings.ToLower(strings.ReplaceAll(name, " ", "_")))
}

func DisplayToCanonical(name string) string {
	return removeColonPerhaps(
		strings.ToLower(strings.ReplaceAll(name, " ", "_")))
}

func CanonicalToDisplay(name string) (res string) {
	tmp := strings.Title(name)
	var afterPoint bool
	for _, ch := range tmp {
		if afterPoint {
			afterPoint = false
			ch = unicode.ToLower(ch)
		}
		switch ch {
		case '.':
			afterPoint = true
		case '_':
			ch = ' '
		}
		res += string(ch)
	}
	return addColonPerhaps(res)
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
