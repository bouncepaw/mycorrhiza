package shroom

import (
	"errors"

	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/markup"
	"github.com/bouncepaw/mycorrhiza/util"
)

func init() {
	markup.HyphaExists = func(hyphaName string) bool {
		return hyphae.ByName(hyphaName).Exists
	}
	markup.HyphaAccess = func(hyphaName string) (rawText, binaryBlock string, err error) {
		if h := hyphae.ByName(hyphaName); h.Exists {
			rawText, err = FetchTextPart(h)
			if h.BinaryPath != "" {
				binaryBlock = BinaryHtmlBlock(h)
			}
		} else {
			err = errors.New("Hypha " + hyphaName + " does not exist")
		}
		return
	}
	markup.HyphaIterate = func(λ func(string)) {
		for h := range hyphae.YieldExistingHyphae() {
			λ(h.Name)
		}
	}
	markup.HyphaImageForOG = func(hyphaName string) string {
		if h := hyphae.ByName(hyphaName); h.Exists && h.BinaryPath != "" {
			return util.URL + "/binary/" + hyphaName
		}
		return util.URL + "/favicon.ico"
	}
}
