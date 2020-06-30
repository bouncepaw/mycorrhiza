package util

import (
	"strings"
	"unicode"
)

func UrlToCanonical(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "_"))
}

func DisplayToCanonical(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "_"))
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
	return res
}
