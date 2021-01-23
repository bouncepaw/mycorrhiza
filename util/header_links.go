package util

import (
	"strings"
)

func SetDefaultHeaderLinks() {
	HeaderLinks = []HeaderLink{
		{"/", SiteName},
		{"/recent-changes", "Recent changes"},
		{"/list", "All hyphae"},
		{"/random", "Random"},
	}
}

// rocketlinkλ is markup.Rocketlink. You have to pass it like that to avoid cyclical dependency.
func ParseHeaderLinks(text string, rocketlinkλ func(string, string) (string, string, string, string)) {
	HeaderLinks = []HeaderLink{}
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "=>") {
			href, text, _, _ := rocketlinkλ(line, HeaderLinksHypha)
			HeaderLinks = append(HeaderLinks, HeaderLink{
				Href:    href,
				Display: text,
			})
		}
	}
}

type HeaderLink struct {
	Href    string
	Display string
}

var HeaderLinks []HeaderLink
