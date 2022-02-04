package shroom

import (
	"errors"

	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/views"

	"github.com/bouncepaw/mycomarkup/v3/globals"
)

func init() {
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
		default:
			rawText, err = FetchTextPart(h)
			if h := h.(*hyphae.MediaHypha); h.Kind() == hyphae.HyphaMedia {
				// the view is localized, but we can't pass it, so...
				binaryBlock = views.AttachmentHTMLRaw(h)
			}
		}
		return
	}
	globals.HyphaIterate = func(λ func(string)) {
		for h := range hyphae.YieldExistingHyphae() {
			λ(h.CanonicalName())
		}
	}
}
