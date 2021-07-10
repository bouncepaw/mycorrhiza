package web

// stuff.go is used for meta stuff about the wiki or all hyphae at once.
import (
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/files"

	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/shroom"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"
)

func initStuff() {
	http.HandleFunc("/list/", handlerList)
	http.HandleFunc("/reindex/", handlerReindex)
	http.HandleFunc("/update-header-links/", handlerUpdateHeaderLinks)
	http.HandleFunc("/random/", handlerRandom)
	http.HandleFunc("/about/", handlerAbout)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, rq *http.Request) {
		http.Redirect(w, rq, "/static/favicon.ico", http.StatusSeeOther)
	})
}

// handlerList shows a list of all hyphae in the wiki in random order.
func handlerList(w http.ResponseWriter, rq *http.Request) {
	u := user.FromRequest(rq)
	if shown := u.ShowLockMaybe(w, rq); shown {
		return
	}
	util.PrepareRq(rq)
	util.HTTP200Page(w, views.BaseHTML("List of pages", views.HyphaListHTML(), u))
}

// handlerReindex reindexes all hyphae by checking the wiki storage directory anew.
func handlerReindex(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	if shown := user.FromRequest(rq).ShowLockMaybe(w, rq); shown {
		return
	}
	if ok := user.CanProceed(rq, "reindex"); !ok {
		httpErr(w, http.StatusForbidden, cfg.HomeHypha, "Not enough rights", "You must be an admin to reindex hyphae.")
		log.Println("Rejected", rq.URL)
		return
	}
	hyphae.ResetCount()
	log.Println("Reindexing hyphae in", files.HyphaeDir())
	hyphae.Index(files.HyphaeDir())
	log.Println("Indexed", hyphae.Count(), "hyphae")
	http.Redirect(w, rq, "/", http.StatusSeeOther)
}

// handlerUpdateHeaderLinks updates header links by reading the configured hypha, if there is any, or resorting to default values.
//
// See https://mycorrhiza.wiki/hypha/configuration/header
func handlerUpdateHeaderLinks(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	if shown := user.FromRequest(rq).ShowLockMaybe(w, rq); shown {
		return
	}
	if ok := user.CanProceed(rq, "update-header-links"); !ok {
		httpErr(w, http.StatusForbidden, cfg.HomeHypha, "Not enough rights", "You must be a moderator to update header links.")
		log.Println("Rejected", rq.URL)
		return
	}
	shroom.SetHeaderLinks()
	http.Redirect(w, rq, "/", http.StatusSeeOther)
}

// handlerRandom redirects to a random hypha.
func handlerRandom(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	if shown := user.FromRequest(rq).ShowLockMaybe(w, rq); shown {
		return
	}
	var (
		randomHyphaName string
		amountOfHyphae  = hyphae.Count()
	)
	if amountOfHyphae == 0 {
		httpErr(w, http.StatusNotFound, cfg.HomeHypha, "There are no hyphae",
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
	if shown := user.FromRequest(rq).ShowLockMaybe(w, rq); shown {
		return
	}
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := io.WriteString(w, views.BaseHTML("About "+cfg.WikiName, views.AboutHTML(), user.FromRequest(rq)))
	if err != nil {
		log.Println(err)
	}
}
