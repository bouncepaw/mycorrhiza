package shroom

import (
	"sort"
	"strings"

	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/util"
)

func YieldHyphaNamesContainingString(query string) <-chan string {
	query = util.CanonicalName(query)
	out := make(chan string)
	go func() {
		// To make it unicode-friendly and lean, we cast every string into rune slices, sort, and only then cast them back
		raw := make([][]rune, 0)
		for h := range hyphae.YieldExistingHyphae() {
			if hyphaNameMatchesString(h.Name, query) {
				raw = append(raw, []rune(h.Name))
			}
		}
		sort.Slice(raw, func(i, j int) bool {
			const slash rune = 47 // == '/'
			// Classic lexicographical sort with a twist
			c := 0
			for {
				if c == len(raw[i]) { return true }
				if c == len(raw[j]) { return false }
				if raw[i][c] == raw[j][c] {
					c++
				} else {
					// The twist: subhyphae-awareness is about pushing slash upwards
					if raw[i][c] == slash { return true }
					if raw[j][c] == slash { return false }
					return raw[i][c] < raw[j][c]
				}
			}
		})
		for _, name := range raw {
			out <- string(name)
		}
		close(out)
	}()
	return out
}

// This thing gotta be changed one day, when a hero has time to implement a good searching algorithm.
func hyphaNameMatchesString(hyphaName, query string) bool {
	return strings.Contains(hyphaName, query)
}
