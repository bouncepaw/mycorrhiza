package interwiki

import (
	"fmt"
	"github.com/bouncepaw/mycorrhiza/util"
	"log"
)

// WikiEngine is an enumeration of supported interwiki targets.
type WikiEngine string

const (
	Mycorrhiza WikiEngine = "mycorrhiza"
	Agora      WikiEngine = "agora"
	// Generic is any website.
	Generic WikiEngine = "generic"
)

func (we WikiEngine) Valid() bool {
	switch we {
	case Mycorrhiza, Agora, Generic:
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

func (w *Wiki) canonize() {
	switch {
	case w.Name == "":
		log.Fatalln("Cannot have a wiki in the interwiki map with no name")
	case w.URL == "":
		log.Fatalf("Wiki ‘%s’ has no URL\n", w.Name)
	case !w.Engine.Valid():
		log.Fatalf("Unknown engine ‘%s’ for wiki ‘%s’\n", w.Engine, w.Name)
	}

	w.Name = util.CanonicalName(w.Name)
	for i, alias := range w.Aliases {
		w.Aliases[i] = util.CanonicalName(alias)
	}

	if w.LinkHrefFormat == "" {
		switch w.Engine {
		case Mycorrhiza:
			w.LinkHrefFormat = fmt.Sprintf("%s/hypha/{NAME}", w.URL)
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
}
