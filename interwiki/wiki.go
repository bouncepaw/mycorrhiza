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

// Wiki is an entry in the interwiki map.
type Wiki struct {
	// Names is a slice of link prefices that correspond to this wiki.
	Names []string `json:"names"`

	// URL is the address of the wiki.
	URL string `json:"url"`

	// LinkFormat is a format string for incoming interwiki links. The format strings should look like this:
	//     http://wiki.example.org/view/%s
	// where %s is where text will be inserted. No other % instructions are supported yet. They will be added once we learn of their use cases.
	//
	// This field is optional. For other wikis, it is automatically set to <URL>/%s; for Mycorrhiza wikis, it is automatically set to <URL>/hypha/%s.
	LinkFormat string `json:"link_format"`

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

	if w.LinkFormat == "" {
		switch w.Engine {
		case Mycorrhiza:
			w.LinkFormat = fmt.Sprintf("%s/hypha/%%s", w.URL)
		default:
			w.LinkFormat = fmt.Sprintf("%s/%%s", w.URL)
		}
	}
}
