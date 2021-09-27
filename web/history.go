package web

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"
)

func initHistory(r *mux.Router) {
	r.PathPrefix("/history/").HandlerFunc(handlerHistory)

	r.HandleFunc("/recent-changes/{count:[0-9]+}", handlerRecentChanges)
	r.HandleFunc("/recent-changes/", func(w http.ResponseWriter, rq *http.Request) {
		http.Redirect(w, rq, "/recent-changes/20", http.StatusSeeOther)
	})

	r.HandleFunc("/recent-changes-rss", handlerRecentChangesRSS)
	r.HandleFunc("/recent-changes-atom", handlerRecentChangesAtom)
	r.HandleFunc("/recent-changes-json", handlerRecentChangesJSON)
}

// handlerHistory lists all revisions of a hypha.
func handlerHistory(w http.ResponseWriter, rq *http.Request) {
	hyphaName := util.HyphaNameFromRq(rq, "history")
	var list string

	// History can be found for files that do not exist anymore.
	revs, err := history.Revisions(hyphaName)
	if err == nil {
		list = history.HistoryWithRevisions(hyphaName, revs)
	}
	log.Println("Found", len(revs), "revisions for", hyphaName)

	var lc = l18n.FromRequest(rq)
	util.HTTP200Page(w, views.BaseHTML(
		fmt.Sprintf(lc.Get("ui.history_title"), util.BeautifulName(hyphaName)),
		views.HistoryHTML(rq, hyphaName, list, lc),
		lc,
		user.FromRequest(rq)))
}

// handlerRecentChanges displays the /recent-changes/ page.
func handlerRecentChanges(w http.ResponseWriter, rq *http.Request) {
	// Error ignored: filtered by regex
	n, _ := strconv.Atoi(mux.Vars(rq)["count"])
	var lc = l18n.FromRequest(rq)
	util.HTTP200Page(w, views.BaseHTML(
		lc.GetPlural("ui.recent_title", n),
		views.RecentChangesHTML(n, lc), 
		lc,
		user.FromRequest(rq)))
}

// genericHandlerOfFeeds is a helper function for the web feed handlers.
func genericHandlerOfFeeds(w http.ResponseWriter, rq *http.Request, f func() (string, error), name string) {
	if content, err := f(); err != nil {
		w.Header().Set("Content-Type", "text/plain;charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "An error while generating "+name+": "+err.Error())
	} else {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, content)
	}
}

func handlerRecentChangesRSS(w http.ResponseWriter, rq *http.Request) {
	genericHandlerOfFeeds(w, rq, history.RecentChangesRSS, "RSS")
}

func handlerRecentChangesAtom(w http.ResponseWriter, rq *http.Request) {
	genericHandlerOfFeeds(w, rq, history.RecentChangesAtom, "Atom")
}

func handlerRecentChangesJSON(w http.ResponseWriter, rq *http.Request) {
	genericHandlerOfFeeds(w, rq, history.RecentChangesJSON, "JSON feed")
}
