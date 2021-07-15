package web

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"
)

func initHistory(r *mux.Router) {
	r.PathPrefix("/history/").HandlerFunc(handlerHistory)

	r.PathPrefix("/recent-changes/").HandlerFunc(handlerRecentChanges)
	r.HandleFunc("/recent-changes-rss", handlerRecentChangesRSS)
	r.HandleFunc("/recent-changes-atom", handlerRecentChangesAtom)
	r.HandleFunc("/recent-changes-json", handlerRecentChangesJSON)
}

// handlerHistory lists all revisions of a hypha.
func handlerHistory(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	if shown := user.FromRequest(rq).ShowLockMaybe(w, rq); shown {
		return
	}
	hyphaName := util.HyphaNameFromRq(rq, "history")
	var list string

	// History can be found for files that do not exist anymore.
	revs, err := history.Revisions(hyphaName)
	if err == nil {
		list = history.HistoryWithRevisions(hyphaName, revs)
	}
	log.Println("Found", len(revs), "revisions for", hyphaName)

	util.HTTP200Page(w,
		views.BaseHTML(hyphaName, views.HistoryHTML(rq, hyphaName, list), user.FromRequest(rq)))
}

// handlerRecentChanges displays the /recent-changes/ page.
func handlerRecentChanges(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	if shown := user.FromRequest(rq).ShowLockMaybe(w, rq); shown {
		return
	}
	var (
		noPrefix = strings.TrimPrefix(rq.URL.String(), "/recent-changes/")
		n, err   = strconv.Atoi(noPrefix)
	)
	if err == nil && n < 101 {
		util.HTTP200Page(w, views.BaseHTML(strconv.Itoa(n)+" recent changes", views.RecentChangesHTML(n), user.FromRequest(rq)))
	} else {
		http.Redirect(w, rq, "/recent-changes/20", http.StatusSeeOther)
	}
}

// genericHandlerOfFeeds is a helper function for the web feed handlers.
func genericHandlerOfFeeds(w http.ResponseWriter, rq *http.Request, f func() (string, error), name string) {
	util.PrepareRq(rq)
	if shown := user.FromRequest(rq).ShowLockMaybe(w, rq); shown {
		return
	}
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
