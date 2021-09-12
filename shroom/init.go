package shroom

import (
	"errors"

	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/views"

	"github.com/bouncepaw/mycomarkup/v2/globals"
)

func init() {
	globals.HyphaExists = func(hyphaName string) bool {
		return hyphae.ByName(hyphaName).Exists
	}
	globals.HyphaAccess = func(hyphaName string) (rawText, binaryBlock string, err error) {
		if h := hyphae.ByName(hyphaName); h.Exists {
			rawText, err = FetchTextPart(h)
			if h.BinaryPath != "" {
				binaryBlock = views.AttachmentHTML(h)
			}
		} else {
			err = errors.New("Hypha " + hyphaName + " does not exist")
		}
		return
	}
	globals.HyphaIterate = func(λ func(string)) {
		for h := range hyphae.YieldExistingHyphae() {
			λ(h.Name)
		}
	}
}
