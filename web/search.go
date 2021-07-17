package web

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/bouncepaw/mycorrhiza/shroom"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"
)

func initSearch(r *mux.Router) {
	r.HandleFunc("/title-search/", handlerTitleSearch)
	r.HandleFunc("/title-search-json/", handlerTitleSearchJSON)
}

func handlerTitleSearch(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	_ = rq.ParseForm()
	var (
		query = rq.FormValue("q")
		u     = user.FromRequest(rq)
	)
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(
		w,
		views.BaseHTML(
			"Search: "+query,
			views.TitleSearchHTML(query, shroom.YieldHyphaNamesContainingString),
			u,
		),
	)
}

func handlerTitleSearchJSON(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	_ = rq.ParseForm()
	query := rq.FormValue("q")
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(
		w,
		views.TitleSearchJSON(query, shroom.YieldHyphaNamesContainingString),
	)
}
