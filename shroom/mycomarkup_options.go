package shroom

import (
	"errors"
	"github.com/bouncepaw/mycomarkup/v5/options"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/interwiki"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"
)

func MarkupOptions(hyphaName string) options.Options {
	return options.Options{
		HyphaName:             hyphaName,
		WebSiteURL:            cfg.URL,
		TransclusionSupported: true,
		RedLinksSupported:     true,
		InterwikiSupported:    true,
		HyphaExists: func(hyphaName string) bool {
			switch hyphae.ByName(hyphaName).(type) {
			case *hyphae.EmptyHypha:
				return false
			default:
				return true
			}
		},
		IterateHyphaNamesWith: func(λ func(string)) {
			for h := range hyphae.YieldExistingHyphae() {
				λ(h.CanonicalName())
			}
		},
		HyphaHTMLData: func(hyphaName string) (rawText, binaryBlock string, err error) {
			switch h := hyphae.ByName(hyphaName).(type) {
			case *hyphae.EmptyHypha:
				err = errors.New("Hypha " + hyphaName + " does not exist")
			case *hyphae.TextualHypha:
				rawText, err = hyphae.FetchMycomarkupFile(h)
			case *hyphae.MediaHypha:
				rawText, err = hyphae.FetchMycomarkupFile(h)
				binaryBlock = views.MediaRaw(h)
			}
			return
		},
		LocalTargetCanonicalName: util.CanonicalName,
		LocalLinkHref: func(hyphaName string) string {
			return "/hypha/" + util.CanonicalName(hyphaName)
		},
		LocalImgSrc: func(hyphaName string) string {
			return "/binary/" + util.CanonicalName(hyphaName)
		},
		LinkHrefFormatForInterwikiPrefix: interwiki.HrefLinkFormatFor,
		ImgSrcFormatForInterwikiPrefix:   interwiki.ImgSrcFormatFor,
	}.FillTheRest()
}
