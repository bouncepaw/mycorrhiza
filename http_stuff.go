// http_stuff.go is used for meta stuff about the wiki or all hyphae at once.
package main

import (
	"github.com/bouncepaw/mycorrhiza/cfg"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/shroom"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"
)

func init() {
	http.HandleFunc("/list/", handlerList)
	http.HandleFunc("/reindex/", handlerReindex)
	http.HandleFunc("/update-header-links/", handlerUpdateHeaderLinks)
	http.HandleFunc("/random/", handlerRandom)
	http.HandleFunc("/about/", handlerAbout)
}

// handlerList shows a list of all hyphae in the wiki in random order.
func handlerList(w http.ResponseWriter, rq *http.Request) {
	prepareRq(rq)
	util.HTTP200Page(w, base("List of pages", views.HyphaListHTML(), user.FromRequest(rq)))
}

// handlerReindex reindexes all hyphae by checking the wiki storage directory anew.
func handlerReindex(w http.ResponseWriter, rq *http.Request) {
	prepareRq(rq)
	if ok := user.CanProceed(rq, "reindex"); !ok {
		HttpErr(w, http.StatusForbidden, cfg.HomeHypha, "Not enough rights", "You must be an admin to reindex hyphae.")
		log.Println("Rejected", rq.URL)
		return
	}
	hyphae.ResetCount()
	log.Println("Wiki storage directory is", WikiDir)
	log.Println("Start indexing hyphae...")
	hyphae.Index(WikiDir)
	log.Println("Indexed", hyphae.Count(), "hyphae")
	http.Redirect(w, rq, "/", http.StatusSeeOther)
}

// handlerUpdateHeaderLinks updates header links by reading the configured hypha, if there is any, or resorting to default values.
//
// See https://mycorrhiza.lesarbr.es/hypha/configuration/header
func handlerUpdateHeaderLinks(w http.ResponseWriter, rq *http.Request) {
	prepareRq(rq)
	if ok := user.CanProceed(rq, "update-header-links"); !ok {
		HttpErr(w, http.StatusForbidden, cfg.HomeHypha, "Not enough rights", "You must be a moderator to update header links.")
		log.Println("Rejected", rq.URL)
		return
	}
	shroom.SetHeaderLinks()
	http.Redirect(w, rq, "/", http.StatusSeeOther)
}

// handlerRandom redirects to a random hypha.
func handlerRandom(w http.ResponseWriter, rq *http.Request) {
	prepareRq(rq)
	var (
		randomHyphaName string
		amountOfHyphae  = hyphae.Count()
	)
	if amountOfHyphae == 0 {
		HttpErr(w, http.StatusNotFound, cfg.HomeHypha, "There are no hyphae",
			"It is impossible to display a random hypha because the wiki does not contain any hyphae")
		return
	}
	i := rand.Intn(amountOfHyphae)
	for h := range hyphae.YieldExistingHyphae() {
		if i == 0 {
			randomHyphaName = h.Name
		}
		i--
	}
	http.Redirect(w, rq, "/hypha/"+randomHyphaName, http.StatusSeeOther)
}

// handlerAbout shows a summary of wiki's software.
func handlerAbout(w http.ResponseWriter, rq *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := io.WriteString(w, base("About "+cfg.WikiName, views.AboutHTML(), user.FromRequest(rq)))
	if err != nil {
		log.Println(err)
	}
}
