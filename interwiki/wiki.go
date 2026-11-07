package interwikisvc

import (
	"fmt"
	"log/slog"

	"github.com/bouncepaw/mycorrhiza/interwiki/domain"
	"github.com/bouncepaw/mycorrhiza/util"
)

// Wiki is an entry in the interwiki map.
type Wiki struct {
	interwikidomain.WikiDTO
}

func (w *Wiki) canonize() error {
	switch {
	case w.Name == "":
		slog.Error("A site in the interwiki map has no name")
		return fmt.Errorf("site with no name")
	case w.URL == "":
		slog.Error("Site in the interwiki map has no URL", "name", w.Name)
		return fmt.Errorf("site with no URL")
	case !w.Engine.Valid():
		slog.Error("Site in the interwiki map has an unknown engine",
			"siteName", w.Name,
			"engine", w.Engine,
		)
		return fmt.Errorf("unknown engine")
	}

	w.Name = util.CanonicalName(w.Name)
	for i, alias := range w.Aliases {
		w.Aliases[i] = util.CanonicalName(alias)
	}

	if w.LinkHrefFormat == "" {
		switch w.Engine {
		case interwikidomain.Mycorrhiza:
			w.LinkHrefFormat = fmt.Sprintf("%s/hypha/{NAME}", w.URL)
		case interwikidomain.Betula:
			w.LinkHrefFormat = fmt.Sprintf("%s/{BETULA-NAME}", w.URL)
		case interwikidomain.Agora:
			w.LinkHrefFormat = fmt.Sprintf("%s/node/{NAME}", w.URL)
		default:
			w.LinkHrefFormat = fmt.Sprintf("%s/{NAME}", w.URL)
		}
	}

	if w.ImgSrcFormat == "" {
		switch w.Engine {
		case interwikidomain.Mycorrhiza:
			w.ImgSrcFormat = fmt.Sprintf("%s/binary/{NAME}", w.URL)
		default:
			w.ImgSrcFormat = fmt.Sprintf("%s/{NAME}", w.URL)
		}
	}

	return nil
}
