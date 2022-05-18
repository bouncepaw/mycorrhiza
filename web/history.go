package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/bouncepaw/mycorrhiza/history"
)

func initHistory(r *mux.Router) {

	r.HandleFunc("/recent-changes-rss", handlerRecentChangesRSS)
	r.HandleFunc("/recent-changes-atom", handlerRecentChangesAtom)
	r.HandleFunc("/recent-changes-json", handlerRecentChangesJSON)
}

// genericHandlerOfFeeds is a helper function for the web feed handlers.
func genericHandlerOfFeeds(w http.ResponseWriter, rq *http.Request, f func(history.FeedOptions) (string, error), name string, contentType string) {
	opts, err := history.ParseFeedOptions(rq.URL.Query())
	var content string
	if err == nil {
		content, err = f(opts)
	}

	if err != nil {
		w.Header().Set("Content-Type", "text/plain;charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "An error while generating "+name+": "+err.Error())
	} else {
		w.Header().Set("Content-Type", fmt.Sprintf("%s;charset=utf-8", contentType))
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, content)
	}
}

func handlerRecentChangesRSS(w http.ResponseWriter, rq *http.Request) {
	genericHandlerOfFeeds(w, rq, history.RecentChangesRSS, "RSS", "application/rss+xml")
}

func handlerRecentChangesAtom(w http.ResponseWriter, rq *http.Request) {
	genericHandlerOfFeeds(w, rq, history.RecentChangesAtom, "Atom", "application/atom+xml")
}

func handlerRecentChangesJSON(w http.ResponseWriter, rq *http.Request) {
	genericHandlerOfFeeds(w, rq, history.RecentChangesJSON, "JSON feed", "application/json")
}
