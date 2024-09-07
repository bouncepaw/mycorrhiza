package shroom

import (
	"strings"

	"github.com/bouncepaw/mycorrhiza/internal/hyphae"
	"github.com/bouncepaw/mycorrhiza/util"
)

// YieldHyphaNamesContainingString picks hyphae with have a string in their title, sorts and iterates over them in alphabetical order.
func YieldHyphaNamesContainingString(query string) <-chan string {
	query = util.CanonicalName(strings.TrimSpace(query))
	out := make(chan string)
	sorted := hyphae.PathographicSort(out)
	go func() {
		for h := range hyphae.YieldExistingHyphae() {
			if hyphaNameMatchesString(h.CanonicalName(), query) {
				out <- h.CanonicalName()
			}
		}
		close(out)
	}()
	return sorted
}

// This thing gotta be changed one day, when a hero has time to implement a good searching algorithm.
func hyphaNameMatchesString(hyphaName, query string) bool {
	return strings.Contains(hyphaName, query)
}
