package web

import (
	"github.com/bouncepaw/mycorrhiza/hyphae/backlinks"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"
)

func initBacklinks(r *mux.Router) {
	r.PathPrefix("/backlinks/").HandlerFunc(handlerBacklinks)
}

// handlerBacklinks lists all backlinks to a hypha.
func handlerBacklinks(w http.ResponseWriter, rq *http.Request) {
	var (
		hyphaName = util.HyphaNameFromRq(rq, "backlinks")
		lc        = l18n.FromRequest(rq)
	)
	util.HTTP200Page(w, views.Base(
		lc.Get("ui.backlinks_title", &l18n.Replacements{"query": util.BeautifulName(hyphaName)}),
		views.Backlinks(hyphaName, backlinks.YieldHyphaBacklinks, lc),
		lc,
		user.FromRequest(rq)))
}
