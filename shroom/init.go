package shroom

import (
	"errors"

	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/views"

	"github.com/bouncepaw/mycomarkup/v3/globals"
)

func init() {
	globals.HyphaExists = func(hyphaName string) bool {
		return hyphae.ByName(hyphaName).DoesExist()
	}
	globals.HyphaAccess = func(hyphaName string) (rawText, binaryBlock string, err error) {
		if h := hyphae.ByName(hyphaName); h.DoesExist() {
			rawText, err = FetchTextPart(h)
			if h := h.(*hyphae.MediaHypha); h.Kind() == hyphae.HyphaMedia {
				// the view is localized, but we can't pass it, so...
				binaryBlock = views.AttachmentHTMLRaw(h)
			}
		} else {
			err = errors.New("MediaHypha " + hyphaName + " does not exist")
		}
		return
	}
	globals.HyphaIterate = func(λ func(string)) {
		for h := range hyphae.YieldExistingHyphae() {
			λ(h.CanonicalName())
		}
	}
}
