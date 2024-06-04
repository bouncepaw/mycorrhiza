package mycoopts

import (
	"errors"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/options"
	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	hyphae2 "github.com/bouncepaw/mycorrhiza/internal/hyphae"
	"github.com/bouncepaw/mycorrhiza/interwiki"
	"github.com/bouncepaw/mycorrhiza/util"
)

func MarkupOptions(hyphaName string) options.Options {
	return options.Options{
		HyphaName:             hyphaName,
		WebSiteURL:            cfg.URL,
		TransclusionSupported: true,
		RedLinksSupported:     true,
		InterwikiSupported:    true,
		HyphaExists: func(hyphaName string) bool {
			switch hyphae2.ByName(hyphaName).(type) {
			case *hyphae2.EmptyHypha:
				return false
			default:
				return true
			}
		},
		IterateHyphaNamesWith: func(λ func(string)) {
			for h := range hyphae2.YieldExistingHyphae() {
				λ(h.CanonicalName())
			}
		},
		HyphaHTMLData: func(hyphaName string) (rawText, binaryBlock string, err error) {
			switch h := hyphae2.ByName(hyphaName).(type) {
			case *hyphae2.EmptyHypha:
				err = errors.New("Hypha " + hyphaName + " does not exist")
			case *hyphae2.TextualHypha:
				rawText, err = hyphae2.FetchMycomarkupFile(h)
			case *hyphae2.MediaHypha:
				rawText, err = hyphae2.FetchMycomarkupFile(h)
				binaryBlock = mediaRaw(h)
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
