package shroom

import (
	"errors"
	"github.com/bouncepaw/mycomarkup/v4/options"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/views"
)

func MarkupOptions(hyphaName string) options.Options {
	return fillMycomarkupOptions(options.Options{
		HyphaName:             hyphaName,
		WebSiteURL:            cfg.URL,
		TransclusionSupported: true,
	})
}

func fillMycomarkupOptions(opts options.Options) options.Options {
	opts.HyphaExists = func(hyphaName string) bool {
		switch hyphae.ByName(hyphaName).(type) {
		case *hyphae.EmptyHypha:
			return false
		default:
			return true
		}
	}
	opts.HyphaHTMLData = func(hyphaName string) (rawText, binaryBlock string, err error) {
		switch h := hyphae.ByName(hyphaName).(type) {
		case *hyphae.EmptyHypha:
			err = errors.New("Hypha " + hyphaName + " does not exist")
		case *hyphae.TextualHypha:
			rawText, err = FetchTextFile(h)
		case *hyphae.MediaHypha:
			rawText, err = FetchTextFile(h)
			binaryBlock = views.MediaRaw(h)
		}
		return
	}
	opts.IterateHyphaNamesWith = func(λ func(string)) {
		for h := range hyphae.YieldExistingHyphae() {
			λ(h.CanonicalName())
		}
	}
	return opts.FillTheRest()
}
