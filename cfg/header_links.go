package cfg

// See https://mycorrhiza.lesarbr.es/hypha/configuration/header
import (
	"strings"
)

// HeaderLinks is a list off current header links. Feel free to iterate it directly but do not modify it by yourself. Call ParseHeaderLinks if you need to set new header links.
var HeaderLinks []HeaderLink

// SetDefaultHeaderLinks sets the header links to the default list of: home hypha, recent changes, hyphae list, random hypha.
func SetDefaultHeaderLinks() {
	HeaderLinks = []HeaderLink{
		{"/", WikiName},
		{"/recent-changes", "Recent changes"},
		{"/list", "All hyphae"},
		{"/random", "Random"},
	}
}

// ParseHeaderLinks extracts all rocketlinks from the given text and saves them as header links. rocketlinkλ must be set to markup.Rocketlink. You have to pass it like that to avoid cyclical dependency.
func ParseHeaderLinks(text string, rocketlinkλ func(string, string) (string, string, string)) {
	HeaderLinks = []HeaderLink{}
	for _, line := range strings.Split(text, "\n") {
		// There is a false positive when parsing markup like that:
		//
		//     ```
		//     => this is not a link, it is part of the preformatted block
		//     ```
		//
		// I do not really care.
		if strings.HasPrefix(line, "=>") {
			href, display, _ := rocketlinkλ(line, HeaderLinksHypha)
			HeaderLinks = append(HeaderLinks, HeaderLink{
				Href:    href,
				Display: display,
			})
		}
	}
}

// HeaderLink represents a header link. Header links are the links shown in the top gray bar.
type HeaderLink struct {
	// Href is the URL of the link. It goes <a href="here">...</a>.
	Href string
	// Display is what is shown when the link is rendered. It goes <a href="...">here</a>.
	Display string
}
