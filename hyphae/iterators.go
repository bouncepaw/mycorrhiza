// File `iterators.go` contains stuff that iterates over hyphae.
package hyphae

import (
	"sort"
	"strings"
)

// YieldExistingHyphae iterates over all hyphae and yields all existing ones.
func YieldExistingHyphae() chan *Hypha {
	ch := make(chan *Hypha)
	go func() {
		for _, h := range byNames {
			if h.Exists {
				ch <- h
			}
		}
		close(ch)
	}()
	return ch
}

// FilterTextHyphae filters the source channel and yields only those hyphae than have text parts.
func FilterTextHyphae(src chan *Hypha) chan *Hypha {
	sink := make(chan *Hypha)
	go func() {
		for h := range src {
			if h.TextPath != "" {
				sink <- h
			}
		}
		close(sink)
	}()
	return sink
}

// PathographicSort sorts paths inside the source channel, preserving the path tree structure
func PathographicSort(src chan string) <-chan string {
	out := make(chan string)
	go func() {
		// To make it unicode-friendly and lean, we cast every string into rune slices, sort, and only then cast them back
		raw := make([][]rune, 0)
		for h := range src {
			raw = append(raw, []rune(h))
		}
		sort.Slice(raw, func(i, j int) bool {
			const slash rune = 47 // == '/'
			// Classic lexicographical sort with a twist
			c := 0
			for {
				if c == len(raw[i]) {
					return true
				}
				if c == len(raw[j]) {
					return false
				}
				if raw[i][c] == raw[j][c] {
					c++
				} else {
					// The twist: subhyphae-awareness is about pushing slash upwards
					if raw[i][c] == slash {
						return true
					}
					if raw[j][c] == slash {
						return false
					}
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

// Subhyphae returns slice of subhyphae.
func (h *Hypha) Subhyphae() []*Hypha {
	hyphae := []*Hypha{}
	for subh := range YieldExistingHyphae() {
		if strings.HasPrefix(subh.Name, h.Name+"/") {
			hyphae = append(hyphae, subh)
		}
	}
	return hyphae
}

// AreFreeNames checks if all given `hyphaNames` are not taken. If they are not taken, `ok` is true. If not, `firstFailure` is the name of the first met hypha that is not free.
func AreFreeNames(hyphaNames ...string) (firstFailure string, ok bool) {
	for h := range YieldExistingHyphae() {
		for _, hn := range hyphaNames {
			if hn == h.Name {
				return hn, false
			}
		}
	}
	return "", true
}
