package interwiki

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/bouncepaw/mycorrhiza/util"
)

// WikiEngine is an enumeration of supported interwiki targets.
type WikiEngine string

const (
	Mycorrhiza WikiEngine = "mycorrhiza"
	Betula     WikiEngine = "betula"
	Agora      WikiEngine = "agora"
	// Generic is any website.
	Generic WikiEngine = "generic"
)

func (we WikiEngine) Valid() bool {
	switch we {
	case Mycorrhiza, Betula, Agora, Generic:
		return true
	}
	return false
}

// Wiki is an entry in the interwiki map.
type Wiki struct {
	// Name is the name of the wiki, and is also one of the possible prefices.
	Name string `json:"name"`

	// Aliases are alternative prefices you can use instead of Name. This slice can be empty.
	Aliases []string `json:"aliases,omitempty"`

	// URL is the address of the wiki.
	URL string `json:"url"`

	// LinkHrefFormat is a format string for interwiki links. See Mycomarkup internal docs hidden deep inside for more information.
	//
	// This field is optional. If it is not set, it is derived from other data. See the code.
	LinkHrefFormat string `json:"link_href_format"`

	ImgSrcFormat string `json:"img_src_format"`

	// Engine is the engine of the wiki. Invalid values will result in a start-up error.
	Engine WikiEngine `json:"engine"`
}

func (w *Wiki) canonize() error {
	switch {
	case w.Name == "":
		slog.Error("A site in the interwiki map has no name")
		return errors.New("site with no name")
	case w.URL == "":
		slog.Error("Site in the interwiki map has no URL", "name", w.Name)
		return errors.New("site with no URL")
	case !w.Engine.Valid():
		slog.Error("Site in the interwiki map has an unknown engine",
			"siteName", w.Name,
			"engine", w.Engine,
		)
		return errors.New("unknown engine")
	}

	w.Name = util.CanonicalName(w.Name)
	for i, alias := range w.Aliases {
		w.Aliases[i] = util.CanonicalName(alias)
	}

	if w.LinkHrefFormat == "" {
		switch w.Engine {
		case Mycorrhiza:
			w.LinkHrefFormat = fmt.Sprintf("%s/hypha/{NAME}", w.URL)
		case Betula:
			w.LinkHrefFormat = fmt.Sprintf("%s/{BETULA-NAME}", w.URL)
		case Agora:
			w.LinkHrefFormat = fmt.Sprintf("%s/node/{NAME}", w.URL)
		default:
			w.LinkHrefFormat = fmt.Sprintf("%s/{NAME}", w.URL)
		}
	}

	if w.ImgSrcFormat == "" {
		switch w.Engine {
		case Mycorrhiza:
			w.ImgSrcFormat = fmt.Sprintf("%s/binary/{NAME}", w.URL)
		default:
			w.ImgSrcFormat = fmt.Sprintf("%s/{NAME}", w.URL)
		}
	}

	return nil
}
