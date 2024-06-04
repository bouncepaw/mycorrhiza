// Package misc provides miscellaneous informative views.
package misc

import (
	"github.com/bouncepaw/mycorrhiza/internal/backlinks"
	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	hyphae2 "github.com/bouncepaw/mycorrhiza/internal/hyphae"
	shroom2 "github.com/bouncepaw/mycorrhiza/internal/shroom"
	"github.com/bouncepaw/mycorrhiza/internal/user"
	"github.com/bouncepaw/mycorrhiza/web/static"
	viewutil2 "github.com/bouncepaw/mycorrhiza/web/viewutil"
	"io"
	"log"
	"math/rand"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"

	"github.com/bouncepaw/mycorrhiza/internal/files"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/util"
)

func InitAssetHandlers(rtr *mux.Router) {
	rtr.HandleFunc("/static/style.css", handlerStyle)
	rtr.HandleFunc("/robots.txt", handlerRobotsTxt)
	rtr.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))))
	rtr.HandleFunc("/favicon.ico", func(w http.ResponseWriter, rq *http.Request) {
		http.Redirect(w, rq, "/static/favicon.ico", http.StatusSeeOther)
	})
}

func InitHandlers(rtr *mux.Router) {
	rtr.HandleFunc("/list", handlerList)
	rtr.HandleFunc("/reindex", handlerReindex)
	rtr.HandleFunc("/update-header-links", handlerUpdateHeaderLinks)
	rtr.HandleFunc("/random", handlerRandom)
	rtr.HandleFunc("/about", handlerAbout)
	rtr.HandleFunc("/title-search/", handlerTitleSearch)
	initViews()
}

// handlerList shows a list of all hyphae in the wiki in random order.
func handlerList(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	// TODO: make this more effective, there are too many loops and vars
	var (
		hyphaNames  = make(chan string)
		sortedHypha = hyphae2.PathographicSort(hyphaNames)
		entries     []listDatum
	)
	for hypha := range hyphae2.YieldExistingHyphae() {
		hyphaNames <- hypha.CanonicalName()
	}
	close(hyphaNames)
	for hyphaName := range sortedHypha {
		switch h := hyphae2.ByName(hyphaName).(type) {
		case *hyphae2.TextualHypha:
			entries = append(entries, listDatum{h.CanonicalName(), ""})
		case *hyphae2.MediaHypha:
			entries = append(entries, listDatum{h.CanonicalName(), filepath.Ext(h.MediaFilePath())[1:]})
		}
	}
	viewList(viewutil2.MetaFrom(w, rq), entries)
}

// handlerReindex reindexes all hyphae by checking the wiki storage directory anew.
func handlerReindex(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	if ok := user.CanProceed(rq, "reindex"); !ok {
		var lc = l18n.FromRequest(rq)
		viewutil2.HttpErr(viewutil2.MetaFrom(w, rq), http.StatusForbidden, cfg.HomeHypha, lc.Get("ui.reindex_no_rights"))
		log.Println("Rejected", rq.URL)
		return
	}
	hyphae2.ResetCount()
	log.Println("Reindexing hyphae in", files.HyphaeDir())
	hyphae2.Index(files.HyphaeDir())
	backlinks.IndexBacklinks()
	http.Redirect(w, rq, "/", http.StatusSeeOther)
}

// handlerUpdateHeaderLinks updates header links by reading the configured hypha, if there is any, or resorting to default values.
func handlerUpdateHeaderLinks(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	if ok := user.CanProceed(rq, "update-header-links"); !ok {
		var lc = l18n.FromRequest(rq)
		viewutil2.HttpErr(viewutil2.MetaFrom(w, rq), http.StatusForbidden, cfg.HomeHypha, lc.Get("ui.header_no_rights"))
		log.Println("Rejected", rq.URL)
		return
	}
	shroom2.SetHeaderLinks()
	http.Redirect(w, rq, "/", http.StatusSeeOther)
}

// handlerRandom redirects to a random hypha.
func handlerRandom(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		randomHyphaName string
		amountOfHyphae  = hyphae2.Count()
	)
	if amountOfHyphae == 0 {
		var lc = l18n.FromRequest(rq)
		viewutil2.HttpErr(viewutil2.MetaFrom(w, rq), http.StatusNotFound, cfg.HomeHypha, lc.Get("ui.random_no_hyphae_tip"))
		return
	}
	i := rand.Intn(amountOfHyphae)
	for h := range hyphae2.YieldExistingHyphae() {
		if i == 0 {
			randomHyphaName = h.CanonicalName()
		}
		i--
	}
	http.Redirect(w, rq, "/hypha/"+randomHyphaName, http.StatusSeeOther)
}

// handlerAbout shows a summary of wiki's software.
func handlerAbout(w http.ResponseWriter, rq *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	var (
		lc    = l18n.FromRequest(rq)
		title = lc.Get("ui.about_title", &l18n.Replacements{"name": cfg.WikiName})
	)
	_, err := io.WriteString(w, viewutil2.Base(
		viewutil2.MetaFrom(w, rq),
		title,
		AboutHTML(lc),
		map[string]string{},
	))
	if err != nil {
		log.Println(err)
	}
}

var stylesheets = []string{"default.css", "custom.css"}

func handlerStyle(w http.ResponseWriter, rq *http.Request) {
	w.Header().Set("Content-Type", mime.TypeByExtension(".css"))
	for _, name := range stylesheets {
		file, err := static.FS.Open(name)
		if err != nil {
			continue
		}
		_, err = io.Copy(w, file)
		if err != nil {
			log.Println(err)
		}
		_ = file.Close()
	}
}

func handlerRobotsTxt(w http.ResponseWriter, rq *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	file, err := static.FS.Open("robots.txt")
	if err != nil {
		return
	}
	_, err = io.Copy(w, file)
	if err != nil {
		log.Println()
	}
	_ = file.Close()
}

func handlerTitleSearch(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	_ = rq.ParseForm()
	var (
		query       = rq.FormValue("q")
		hyphaName   = util.CanonicalName(query)
		_, nameFree = hyphae2.AreFreeNames(hyphaName)
		results     []string
	)
	for hyphaName := range shroom2.YieldHyphaNamesContainingString(query) {
		results = append(results, hyphaName)
	}
	w.WriteHeader(http.StatusOK)
	viewTitleSearch(viewutil2.MetaFrom(w, rq), query, hyphaName, !nameFree, results)
}
