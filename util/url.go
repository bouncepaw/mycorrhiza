/* This file implements things defined by Wikilink RFC. See :main/help/wikilink
 */
package util

import (
	"path"
	"regexp"
	"strings"
)

// `name` must be non-empty.
func sections(name string) (mycel, hyphaName string) {
	mycelRe := regexp.MustCompile(`^:.*/`)
	loc := mycelRe.FindIndex([]byte(name))
	if loc != nil { // if has mycel
		mycel = name[:loc[1]]
		name = name[loc[1]:]
	}
	return mycel, name
}

// Wikilink processes `link` as defined by :main/help/wikilink assuming that `atHypha` is current hypha name.
func Wikilink(link, atHypha string) string {
	mycel, hyphaName := sections(atHypha)
	urlProtocolRe := regexp.MustCompile(`^[a-zA-Z]+:`)
	switch {
	case strings.HasPrefix(link, "::"):
		return "/" + mycel + link[2:]
	case strings.HasPrefix(link, ":"):
		return "/" + link
	case strings.HasPrefix(link, "../") && strings.Count(hyphaName, "/") > 0:
		return "/" + path.Dir(atHypha) + "/" + link[3:]
	case strings.HasPrefix(link, "../"):
		return "/" + mycel + link[3:]
	case strings.HasPrefix(link, "/"):
		return "/" + atHypha + link
	case strings.HasPrefix(link, "./"):
		return "/" + atHypha + link[1:]
	case urlProtocolRe.MatchString(link):
		return link
	default:
		return "/" + link
	}
}
