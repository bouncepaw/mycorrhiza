package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/bouncepaw/mycorrhiza/history"
	"github.com/bouncepaw/mycorrhiza/templates"
	"github.com/bouncepaw/mycorrhiza/util"
)

func init() {
	http.HandleFunc("/history/", handlerHistory)
	http.HandleFunc("/recent-changes/", handlerRecentChanges)
	http.HandleFunc("/recent-changes-rss", handlerRecentChangesRSS)
	http.HandleFunc("/recent-changes-atom", handlerRecentChangesAtom)
	http.HandleFunc("/recent-changes-json", handlerRecentChangesJSON)
}

// handlerHistory lists all revisions of a hypha
func handlerHistory(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	hyphaName := HyphaNameFromRq(rq, "history")
	var list string

	// History can be found for files that do not exist anymore.
	revs, err := history.Revisions(hyphaName)
	if err == nil {
		list = history.HistoryWithRevisions(hyphaName, revs)
	}
	log.Println("Found", len(revs), "revisions for", hyphaName)

	util.HTTP200Page(w,
		base(hyphaName, templates.HistoryHTML(rq, hyphaName, list)))
}

// Recent changes
func handlerRecentChanges(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		noPrefix = strings.TrimPrefix(rq.URL.String(), "/recent-changes/")
		n, err   = strconv.Atoi(noPrefix)
	)
	if err == nil && n < 101 {
		util.HTTP200Page(w, base(strconv.Itoa(n)+" recent changes", history.RecentChanges(n)))
	} else {
		http.Redirect(w, rq, "/recent-changes/20", http.StatusSeeOther)
	}
}

func genericHandlerOfFeeds(w http.ResponseWriter, rq *http.Request, f func() (string, error), name string) {
	log.Println(rq.URL)
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
