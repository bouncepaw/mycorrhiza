package web

// stuff.go is used for meta stuff about the wiki or all hyphae at once.
import (
	"github.com/bouncepaw/mycorrhiza/hyphae/backlinks"
	"github.com/bouncepaw/mycorrhiza/viewutil"
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
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/shroom"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"

	"github.com/bouncepaw/mycomarkup/v3"
	"github.com/bouncepaw/mycomarkup/v3/mycocontext"
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
	lc := l18n.FromRequest(rq)
	articlePath := strings.TrimPrefix(strings.TrimPrefix(rq.URL.Path, "/help/"), "/help")
	// See the history of this file to resurrect the old algorithm that supported multiple languages
	lang := "en"
	if articlePath == "" {
		articlePath = "en"
	}

	if !strings.HasPrefix(articlePath, "en") {
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, "404 Not found")
		return
	}

	content, err := help.Get(articlePath)
	if err != nil && strings.HasPrefix(err.Error(), "open") {
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(
			w,
			views.Base(
				viewutil.MetaFrom(w, rq),
				lc.Get("help.entry_not_found"),
				views.Help(views.HelpEmptyError(lc), lang, lc),
			),
		)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(
			w,
			views.Base(
				viewutil.MetaFrom(w, rq),
				err.Error(),
				views.Help(err.Error(), lang, lc),
			),
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
		views.Base(
			viewutil.MetaFrom(w, rq),
			lc.Get("help.title"),
			views.Help(result, lang, lc),
		),
	)
}

// handlerList shows a list of all hyphae in the wiki in random order.
func handlerList(w http.ResponseWriter, rq *http.Request) {
	var lc = l18n.FromRequest(rq)
	util.PrepareRq(rq)
	util.HTTP200Page(w, views.Base(
		viewutil.MetaFrom(w, rq),
		lc.Get("ui.list_title"),
		views.HyphaList(lc),
	))
}

// handlerReindex reindexes all hyphae by checking the wiki storage directory anew.
func handlerReindex(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	if ok := user.CanProceed(rq, "reindex"); !ok {
		var lc = l18n.FromRequest(rq)
		httpErr(viewutil.MetaFrom(w, rq), lc, http.StatusForbidden, cfg.HomeHypha, lc.Get("ui.reindex_no_rights"))
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
//
// See https://mycorrhiza.wiki/hypha/configuration/header
func handlerUpdateHeaderLinks(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	if ok := user.CanProceed(rq, "update-header-links"); !ok {
		var lc = l18n.FromRequest(rq)
		httpErr(viewutil.MetaFrom(w, rq), lc, http.StatusForbidden, cfg.HomeHypha, lc.Get("ui.header_no_rights"))
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
		httpErr(viewutil.MetaFrom(w, rq), lc, http.StatusNotFound, cfg.HomeHypha, lc.Get("ui.random_no_hyphae_tip"))
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
