package web

import (
	"github.com/bouncepaw/mycorrhiza/hyphae/backlinks"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/bouncepaw/mycorrhiza/l18n"
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
	var (
		hyphaName = util.HyphaNameFromRq(rq, "backlinks")
		lc        = l18n.FromRequest(rq)
	)
	util.HTTP200Page(w, views.BaseHTML(
		lc.Get("ui.backlinks_title", &l18n.Replacements{"query": util.BeautifulName(hyphaName)}),
		views.BacklinksHTML(hyphaName, backlinks.YieldHyphaBacklinks, lc),
		lc,
		user.FromRequest(rq)))
}

func handlerBacklinksJSON(w http.ResponseWriter, rq *http.Request) {
	hyphaName := util.HyphaNameFromRq(rq, "backlinks")
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(
		w,
		views.TitleSearchJSON(hyphaName, backlinks.YieldHyphaBacklinks),
	)
}
