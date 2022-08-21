package shroom

import (
	"git.sr.ht/~bouncepaw/mycomarkup/v5"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/blocks"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/mycocontext"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/mycoopts"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"os"
)

// SetHeaderLinks initializes header links by reading the configured hypha, if there is any, or resorting to default values.
func SetHeaderLinks() {
	switch userLinksHypha := hyphae.ByName(cfg.HeaderLinksHypha).(type) {
	case *hyphae.EmptyHypha:
		setDefaultHeaderLinks()
	case hyphae.ExistingHypha:
		contents, err := os.ReadFile(userLinksHypha.TextFilePath())
		if err != nil || len(contents) == 0 {
			setDefaultHeaderLinks()
		} else {
			text := string(contents)
			parseHeaderLinks(text)
		}
	}
}

// setDefaultHeaderLinks sets the header links to the default list of: home hypha, recent changes, hyphae list, random hypha.
func setDefaultHeaderLinks() {
	viewutil.HeaderLinks = []viewutil.HeaderLink{
		{"/recent-changes", "Recent changes"},
		{"/list", "All hyphae"},
		{"/random", "Random"},
		{"/help", "Help"},
		{"/category", "Categories"},
	}
}

// parseHeaderLinks extracts all rocketlinks from the given text and saves them as header links.
func parseHeaderLinks(text string) {
	viewutil.HeaderLinks = []viewutil.HeaderLink{}
	ctx, _ := mycocontext.ContextFromStringInput(text, mycoopts.MarkupOptions(""))
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
