// Package shroom provides utilities for hypha manipulation.
//
// Some of them are wrappers around functions provided by package hyphae. They manage history for you.
package shroom

import (
	"errors"

	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/views"

	"github.com/bouncepaw/mycomarkup/v3/globals"
)

func init() {
	// TODO: clean this complete and utter mess
	globals.HyphaExists = func(hyphaName string) bool {
		switch hyphae.ByName(hyphaName).(type) {
		case *hyphae.EmptyHypha:
			return false
		default:
			return true
		}
	}
	globals.HyphaAccess = func(hyphaName string) (rawText, binaryBlock string, err error) {
		switch h := hyphae.ByName(hyphaName).(type) {
		case *hyphae.EmptyHypha:
			err = errors.New("Hypha " + hyphaName + " does not exist")
		case *hyphae.TextualHypha:
			rawText, err = FetchTextFile(h)
		case *hyphae.MediaHypha:
			rawText, err = FetchTextFile(h)
			binaryBlock = views.MediaHTMLRaw(h)
		}
		return
	}
	globals.HyphaIterate = func(λ func(string)) {
		for h := range hyphae.YieldExistingHyphae() {
			λ(h.CanonicalName())
		}
	}
}
