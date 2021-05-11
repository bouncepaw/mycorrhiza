package markup

import (
	"strings"

	"github.com/bouncepaw/mycomarkup/links"
)

// LinkParts determines what href, text and class should resulting <a> have based on mycomarkup's addr, display and hypha name.
//
// => addr display
// [[addr|display]]
// TODO: deprecate
func LinkParts(addr, display, hyphaName string) (href, text, class string) {
	l := links.From(addr, display, hyphaName)
	if l.OfKind(links.LinkLocalHypha) && !HyphaExists(l.Address()) {
		l.DestinationUnknown = true
	}
	return l.Href(), l.Display(), l.Classes()
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
