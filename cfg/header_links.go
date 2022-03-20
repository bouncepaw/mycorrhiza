package cfg

// See https://mycorrhiza.wiki/hypha/configuration/header
import (
	"github.com/bouncepaw/mycomarkup/v3"
	"github.com/bouncepaw/mycomarkup/v3/blocks"
	"github.com/bouncepaw/mycomarkup/v3/mycocontext"
)

// HeaderLinks is a list off current header links. Feel free to iterate it directly but do not modify it by yourself. Call ParseHeaderLinks if you need to set new header links.
var HeaderLinks []HeaderLink

// SetDefaultHeaderLinks sets the header links to the default list of: home hypha, recent changes, hyphae list, random hypha.
func SetDefaultHeaderLinks() {
	HeaderLinks = []HeaderLink{
		{"/recent-changes", "Recent changes"},
		{"/list", "All hyphae"},
		{"/random", "Random"},
		{"/help", "Help"},
		{"/category", "Categories"},
	}
}

// ParseHeaderLinks extracts all rocketlinks from the given text and saves them as header links.
func ParseHeaderLinks(text string) {
	HeaderLinks = []HeaderLink{}
	ctx, _ := mycocontext.ContextFromStringInput("", text)
	// We call for side-effects
	_ = mycomarkup.BlockTree(ctx, func(block blocks.Block) {
		switch launchpad := block.(type) {
		case blocks.LaunchPad:
			for _, rocket := range launchpad.Rockets {
				HeaderLinks = append(HeaderLinks, HeaderLink{
					Href:    rocket.Href(),
					Display: rocket.Display(),
				})
			}
		}
	})
}

// HeaderLink represents a header link. Header links are the links shown in the top gray bar.
type HeaderLink struct {
	// Href is the URL of the link. It goes <a href="here">...</a>.
	Href string
	// Display is what is shown when the link is rendered. It goes <a href="...">here</a>.
	Display string
}
