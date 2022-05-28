package interwiki

import (
	"fmt"
	"github.com/bouncepaw/mycorrhiza/util"
	"log"
)

// WikiEngine is an enumeration of supported interwiki targets.
type WikiEngine int

const (
	Mycorrhiza WikiEngine = iota
	OddMuse
	MediaWiki
	MoinMoin1
	MoinMoin2
	DokuWiki
	// Generic is any website.
	Generic
)

// EmojiWithName returns a Unicode emoji that kinda represents the engine and the engine name. One day we might move to actual images. OK for now.
func (we WikiEngine) EmojiWithName() string {
	switch we {
	case Mycorrhiza:
		return "🍄 Mycorrhiza"
	case OddMuse:
		return "🐫 OddMuse"
	case MediaWiki:
		return "🌻 MediaWiki"
	case MoinMoin1:
		return "Ⓜ️ MoinMoin 1.9"
	case MoinMoin2:
		return "Ⓜ️ MoinMoin 2.*"
	case DokuWiki:
		return "📝 DokuWiki"
	default:
		return "🌐 Generic"
	}
}

// Wiki is an entry in the interwiki map.
type Wiki struct {
	// Names is a slice of link prefices that correspond to this wiki.
	Names []string `json:"names"`

	// URL is the address of the wiki.
	URL string `json:"url"`

	// LinkHrefFormat is a format string for interwiki links. See Mycomarkup internal docs hidden deep inside for more information.
	//
	// This field is optional. For other wikis, it is automatically set to <URL>/{NAME}; for Mycorrhiza wikis, it is automatically set to <URL>/hypha/{NAME}}.
	LinkHrefFormat string `json:"link_href_format"`

	ImgSrcFormat string `json:"img_src_format"`

	// Description is a plain-text description of the wiki.
	Description string `json:"description"`

	// Engine is the engine of the wiki. This field is not set in JSON.
	Engine WikiEngine `json:"-"`

	// EngineString is a string name of the engine. It is then converted to Engine. See the code to learn the supported values. All other values will result in an error.
	EngineString string `json:"engine"`
}

var wikiEnginesLookup = map[string]WikiEngine{
	"mycorrhiza": Mycorrhiza,
	"oddmuse":    OddMuse,
	"mediawiki":  MediaWiki,
	"moin1":      MoinMoin1,
	"moin2":      MoinMoin2,
	"dokuwiki":   DokuWiki,
	"generic":    Generic,
}

func (w *Wiki) canonize() {
	if engine, ok := wikiEnginesLookup[w.EngineString]; ok {
		w.Engine = engine
		w.EngineString = "" // Ain't gonna need it anymore
	} else {
		log.Fatalf("Unknown engine ‘%s’\n", w.EngineString)
	}

	if len(w.Names) == 0 {
		log.Fatalln("Cannot have a wiki in the interwiki map with no name")
	}

	if w.URL == "" {
		log.Fatalf("Wiki ‘%s’ has no URL\n", w.Names[0])
	}

	for i, prefix := range w.Names {
		w.Names[i] = util.CanonicalName(prefix)
	}

	if w.LinkHrefFormat == "" {
		switch w.Engine {
		case Mycorrhiza:
			w.LinkHrefFormat = fmt.Sprintf("%s/hypha/{NAME}", w.URL)
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
