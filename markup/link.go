package markup

import (
	"fmt"
	"path"
	"strings"
)

// LinkParts determines what href, text and class should resulting <a> have based on mycomarkup's addr, display and hypha name.
//
// => addr display
// [[addr|display]]
func LinkParts(addr, display, hyphaName string) (href, text, class string) {
	if display == "" {
		text = addr
	} else {
		text = strings.TrimSpace(display)
	}
	class = "wikilink_internal"

	switch {
	case strings.ContainsRune(addr, ':'):
		pos := strings.IndexRune(addr, ':')
		destination := addr[:pos]
		text = addr[pos+1:]
		if strings.HasPrefix(text, "//") && len(text) > 2 {
			text = text[2:]
		}
		return addr, text + fmt.Sprintf(`<img class="wikilink__destination-type" src="/static/icon/%s" width="16" height="16"/>`, destination), "wikilink_external"
	case strings.HasPrefix(addr, "/"):
		return addr, text, class
	case strings.HasPrefix(addr, "./"):
		hyphaName = canonicalName(path.Join(hyphaName, addr[2:]))
	case strings.HasPrefix(addr, "../"):
		hyphaName = canonicalName(path.Join(path.Dir(hyphaName), addr[3:]))
	default:
		hyphaName = canonicalName(addr)
	}
	if !HyphaExists(hyphaName) {
		class += " wikilink_new"
	}
	return "/page/" + hyphaName, text, class
}

// Parse markup line starting with "=>" according to wikilink rules.
// See http://localhost:1737/page/wikilink
func Rocketlink(src, hyphaName string) (href, text, class string) {
	src = strings.TrimSpace(src[2:]) // Drop =>
	if src == "" {
		return
	}
	// Href is text after => till first whitespace
	addr := strings.Fields(src)[0]
	display := strings.TrimPrefix(src, addr)
	return LinkParts(addr, display, hyphaName)
}
