package shroom

import (
	"strings"

	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/util"
)

// YieldHyphaNamesContainingString picks hyphae with have a string in their title, sorts and iterates over them.
func YieldHyphaNamesContainingString(query string) <-chan string {
	query = util.CanonicalName(query)
	out := make(chan string)
	sorted := hyphae.PathographicSort(out)
	go func() {
		for h := range hyphae.YieldExistingHyphae() {
			if hyphaNameMatchesString(h.Name, query) {
				out <- h.Name
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
