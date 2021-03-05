// File `iterators.go` contains stuff that iterates over hyphae.
package hyphae

import (
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
