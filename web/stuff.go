package web

// stuff.go is used for meta stuff about the wiki or all hyphae at once.
import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/files"
	"github.com/bouncepaw/mycorrhiza/help"
	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/shroom"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"

	"github.com/bouncepaw/mycomarkup"
	"github.com/bouncepaw/mycomarkup/mycocontext"
)

func initStuff(r *mux.Router) {
	r.PathPrefix("/help").HandlerFunc(handlerHelp)
	r.HandleFunc("/list", handlerList)
	r.HandleFunc("/reindex", handlerReindex)
	r.HandleFunc("/update-header-links", handlerUpdateHeaderLinks)
	r.HandleFunc("/random", handlerRandom)
	r.HandleFunc("/about", handlerAbout)
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, rq *http.Request) {
		http.Redirect(w, rq, "/static/favicon.ico", http.StatusSeeOther)
	})
}

// handlerHelp gets the appropriate documentation or tells you where you (personally) have failed.
func handlerHelp(w http.ResponseWriter, rq *http.Request) {
	articlePath := strings.TrimPrefix(strings.TrimPrefix(rq.URL.Path, "/help/"), "/help")
	if articlePath == "" {
		articlePath = "en"
	}
	content, err := help.Get(articlePath)
	if err != nil && strings.HasPrefix(err.Error(), "open") {
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(
			w,
			views.BaseHTML("Entry not found",
				views.HelpHTML(views.HelpEmptyErrorHTML()),
				user.FromRequest(rq)),
		)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(
			w,
			views.BaseHTML(err.Error(),
				views.HelpHTML(err.Error()),
				user.FromRequest(rq)),
		)
		return
	}

	// TODO: change for the function that uses byte array when there is such function in mycomarkup.
	ctx, _ := mycocontext.ContextFromStringInput(articlePath, string(content))
	ast := mycomarkup.BlockTree(ctx)
	result := mycomarkup.BlocksToHTML(ctx, ast)
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(
		w,
		views.BaseHTML("Help",
			views.HelpHTML(result),
			user.FromRequest(rq)),
	)
}

// handlerList shows a list of all hyphae in the wiki in random order.
func handlerList(w http.ResponseWriter, rq *http.Request) {
	u := user.FromRequest(rq)
	util.PrepareRq(rq)
	util.HTTP200Page(w, views.BaseHTML("List of pages", views.HyphaListHTML(), u))
}

// handlerReindex reindexes all hyphae by checking the wiki storage directory anew.
func handlerReindex(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	if ok := user.CanProceed(rq, "reindex"); !ok {
		httpErr(w, http.StatusForbidden, cfg.HomeHypha, "Not enough rights", "You must be an admin to reindex hyphae.")
		log.Println("Rejected", rq.URL)
		return
	}
	hyphae.ResetCount()
	log.Println("Reindexing hyphae in", files.HyphaeDir())
	hyphae.Index(files.HyphaeDir())
	http.Redirect(w, rq, "/", http.StatusSeeOther)
}

// handlerUpdateHeaderLinks updates header links by reading the configured hypha, if there is any, or resorting to default values.
//
// See https://mycorrhiza.wiki/hypha/configuration/header
func handlerUpdateHeaderLinks(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
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
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := io.WriteString(w, views.BaseHTML("About "+cfg.WikiName, views.AboutHTML(), user.FromRequest(rq)))
	if err != nil {
		log.Println(err)
	}
}
