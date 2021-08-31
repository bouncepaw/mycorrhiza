package web

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/bouncepaw/mycorrhiza/hyphae"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"
)

func initBacklinks(r *mux.Router) {
	r.PathPrefix("/backlinks/").HandlerFunc(handlerBacklinks)
	r.PathPrefix("/backlinks-json/").HandlerFunc(handlerBacklinksJSON)
}

// handlerBacklinks lists all backlinks to a hypha.
func handlerBacklinks(w http.ResponseWriter, rq *http.Request) {
	hyphaName := util.HyphaNameFromRq(rq, "backlinks")
	util.HTTP200Page(w, views.BaseHTML(
		fmt.Sprintf("Backlinks to %s", util.BeautifulName(hyphaName)),
		views.BacklinksHTML(hyphaName, hyphae.YieldHyphaBacklinks),
		user.FromRequest(rq)))
}

func handlerBacklinksJSON(w http.ResponseWriter, rq *http.Request) {
	hyphaName := util.HyphaNameFromRq(rq, "backlinks")
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(
		w,
		views.TitleSearchJSON(hyphaName, hyphae.YieldHyphaBacklinks),
	)
}
