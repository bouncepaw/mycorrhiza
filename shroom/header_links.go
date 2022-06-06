package shroom

import (
	"github.com/bouncepaw/mycomarkup/v5"
	"github.com/bouncepaw/mycomarkup/v5/blocks"
	"github.com/bouncepaw/mycomarkup/v5/mycocontext"
	"github.com/bouncepaw/mycorrhiza/viewutil"
)

// SetDefaultHeaderLinks sets the header links to the default list of: home hypha, recent changes, hyphae list, random hypha.
func SetDefaultHeaderLinks() {
	viewutil.HeaderLinks = []viewutil.HeaderLink{
		{"/recent-changes", "Recent changes"},
		{"/list", "All hyphae"},
		{"/random", "Random"},
		{"/help", "Help"},
		{"/category", "Categories"},
	}
}

// ParseHeaderLinks extracts all rocketlinks from the given text and saves them as header links.
func ParseHeaderLinks(text string) {
	viewutil.HeaderLinks = []viewutil.HeaderLink{}
	ctx, _ := mycocontext.ContextFromStringInput(text, MarkupOptions(""))
	// We call for side-effects
	_ = mycomarkup.BlockTree(ctx, func(block blocks.Block) {
		switch launchpad := block.(type) {
		case blocks.LaunchPad:
			for _, rocket := range launchpad.Rockets {
				viewutil.HeaderLinks = append(viewutil.HeaderLinks, viewutil.HeaderLink{
					Href:    rocket.LinkHref(ctx),
					Display: rocket.DisplayedText(),
				})
			}
		}
	})
}
