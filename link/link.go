package link

import (
	"fmt"
	"path"
	"strings"

	"github.com/bouncepaw/mycorrhiza/util"
)

// LinkType tells what type the given link is.
type LinkType int

const (
	LinkInavild LinkType = iota
	// LinkLocalRoot is a link like "/list", "/user-list", etc.
	LinkLocalRoot
	// LinkLocalHypha is a link like "test", "../test", etc.
	LinkLocalHypha
	// LinkExternal is an external link with specified protocol.
	LinkExternal
	// LinkInterwiki is currently unused.
	LinkInterwiki
)

// Link is an abstraction for universal representation of links, be they links in mycomarkup links or whatever.
type Link struct {
	// Address is what the link points to.
	Address string
	// Display is what gets nested into the <a> tag.
	Display            string
	Kind               LinkType
	DestinationUnknown bool

	Protocol string
	// How the link address looked originally in source text.
	SrcAddress string
	// How the link display text looked originally in source text. May be empty.
	SrcDisplay string
	// RelativeTo is hypha name to which the link is relative to.
	RelativeTo string
}

// DoubtExistence sets DestinationUnknown to true if the link is local hypha link.
func (l *Link) DoubtExistence() {
	if l.Kind == LinkLocalHypha {
		l.DestinationUnknown = true
	}
}

// Classes returns CSS class string for given link.
func (l *Link) Classes() string {
	if l.Kind == LinkExternal {
		return fmt.Sprintf("wikilink wikilink_external wikilink_%s", l.Protocol)
	}
	classes := "wikilink wikilink_internal"
	if l.DestinationUnknown {
		classes += " wikilink_new"
	}
	return classes
}

// Href returns content for the href attrubite for hyperlink. You should always use it.
func (l *Link) Href() string {
	switch l.Kind {
	case LinkExternal, LinkLocalRoot:
		return l.Address
	default:
		return "/hypha/" + l.Address
	}
}

// ImgSrc returns content for src attribute of img tag. Used with `img{}`.
func (l *Link) ImgSrc() string {
	switch l.Kind {
	case LinkExternal, LinkLocalRoot:
		return l.Address
	default:
		return "/binary/" + l.Address
	}
}

// From returns a Link object given these `address` and `display` on relative to given `hyphaName`.
func From(address, display, hyphaName string) *Link {
	address = strings.TrimSpace(address)
	link := Link{
		SrcAddress: address,
		SrcDisplay: display,
		RelativeTo: hyphaName,
	}

	if display == "" {
		link.Display = address
	} else {
		link.Display = strings.TrimSpace(display)
	}

	switch {
	case strings.ContainsRune(address, ':'):
		pos := strings.IndexRune(address, ':')
		link.Protocol = address[:pos]
		link.Kind = LinkExternal

		if display == "" {
			link.Display = address[pos+1:]
			if strings.HasPrefix(link.Display, "//") && len(link.Display) > 2 {
				link.Display = link.Display[2:]
			}
		}
		link.Address = address
	case strings.HasPrefix(address, "/"):
		link.Address = address
		link.Kind = LinkLocalRoot
	case strings.HasPrefix(address, "./"):
		link.Kind = LinkLocalHypha
		link.Address = util.CanonicalName(path.Join(hyphaName, address[2:]))
	case strings.HasPrefix(address, "../"):
		link.Kind = LinkLocalHypha
		link.Address = util.CanonicalName(path.Join(path.Dir(hyphaName), address[3:]))
	default:
		link.Kind = LinkLocalHypha
		link.Address = util.CanonicalName(address)
	}

	return &link
}
