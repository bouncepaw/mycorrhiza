package interwiki

import (
	"embed"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	//go:embed *html
	fs             embed.FS
	ruTranslation  = ``
	chainInterwiki viewutil.Chain
)

func InitHandlers(rtr *mux.Router) {
	chainInterwiki = viewutil.CopyEnRuWith(fs, "view_interwiki.html", ruTranslation)
	rtr.HandleFunc("/interwiki", handlerInterwiki)
}

func handlerInterwiki(w http.ResponseWriter, rq *http.Request) {
	viewInterwiki(viewutil.MetaFrom(w, rq))
}

type interwikiData struct {
	*viewutil.BaseData
	Entries []*Wiki
	// Emojies contains emojies that represent wiki engines. Emojies[i] is an emoji for Entries[i].Engine
	Emojies []string
	CanEdit bool
}

func viewInterwiki(meta viewutil.Meta) {
	viewutil.ExecutePage(meta, chainInterwiki, interwikiData{
		BaseData: &viewutil.BaseData{},
		Entries:  theMap.list,
		Emojies:  emojiesForEngines(theMap.list),
		CanEdit:  meta.U.Group == "admin",
	})
}

func emojiesForEngines(list []*Wiki) (emojies []string) {
	for _, entry := range list {
		emojies = append(emojies, entry.Engine.EmojiWithName())
	}
	return emojies
}
