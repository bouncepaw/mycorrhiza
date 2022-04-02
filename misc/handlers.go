// Package misc provides miscellaneous informative views.
package misc

import (
	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/files"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/hyphae/backlinks"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/shroom"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"github.com/gorilla/mux"
	"io"
	"log"
	"math/rand"
	"net/http"
)

func InitHandlers(rtr *mux.Router) {
	rtr.HandleFunc("/list", handlerList)
	rtr.HandleFunc("/reindex", handlerReindex)
	rtr.HandleFunc("/update-header-links", handlerUpdateHeaderLinks)
	rtr.HandleFunc("/random", handlerRandom)
	rtr.HandleFunc("/about", handlerAbout)
	rtr.HandleFunc("/favicon.ico", func(w http.ResponseWriter, rq *http.Request) {
		http.Redirect(w, rq, "/static/favicon.ico", http.StatusSeeOther)
	})
	initViews()
}

// handlerList shows a list of all hyphae in the wiki in random order.
func handlerList(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	viewList(viewutil.MetaFrom(w, rq))
}

// handlerReindex reindexes all hyphae by checking the wiki storage directory anew.
func handlerReindex(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	if ok := user.CanProceed(rq, "reindex"); !ok {
		var lc = l18n.FromRequest(rq)
		viewutil.HttpErr(viewutil.MetaFrom(w, rq), lc, http.StatusForbidden, cfg.HomeHypha, lc.Get("ui.reindex_no_rights"))
		log.Println("Rejected", rq.URL)
		return
	}
	hyphae.ResetCount()
	log.Println("Reindexing hyphae in", files.HyphaeDir())
	hyphae.Index(files.HyphaeDir())
	backlinks.IndexBacklinks()
	http.Redirect(w, rq, "/", http.StatusSeeOther)
}

// handlerUpdateHeaderLinks updates header links by reading the configured hypha, if there is any, or resorting to default values.
func handlerUpdateHeaderLinks(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	if ok := user.CanProceed(rq, "update-header-links"); !ok {
		var lc = l18n.FromRequest(rq)
		viewutil.HttpErr(viewutil.MetaFrom(w, rq), lc, http.StatusForbidden, cfg.HomeHypha, lc.Get("ui.header_no_rights"))
		log.Println("Rejected", rq.URL)
		return
	}
	shroom.SetHeaderLinks()
	http.Redirect(w, rq, "/", http.StatusSeeOther)
}

// handlerRandom redirects to a random hypha.
func handlerRandom(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		randomHyphaName string
		amountOfHyphae  = hyphae.Count()
	)
	if amountOfHyphae == 0 {
		var lc = l18n.FromRequest(rq)
		viewutil.HttpErr(viewutil.MetaFrom(w, rq), lc, http.StatusNotFound, cfg.HomeHypha, lc.Get("ui.random_no_hyphae_tip"))
		return
	}
	i := rand.Intn(amountOfHyphae)
	for h := range hyphae.YieldExistingHyphae() {
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
	_, err := io.WriteString(w, views.Base(
		viewutil.MetaFrom(w, rq),
		title,
		views.AboutHTML(lc),
	))
	if err != nil {
		log.Println(err)
	}
}
