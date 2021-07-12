package web

import (
	"github.com/bouncepaw/mycorrhiza/shroom"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/views"
	"io"
	"net/http"

	"github.com/bouncepaw/mycorrhiza/util"
)

func initSearch() {
	http.HandleFunc("/primitive-search/", handlerPrimitiveSearch)
}

func handlerPrimitiveSearch(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		// It just so happened that this function does what we need! Sorry for party rocking.
		query = util.HyphaNameFromRq(rq, "primitive-search")
		u     = user.FromRequest(rq)
	)
	_, _ = io.WriteString(
		w,
		views.BaseHTML(
			"Search: "+query,
			views.PrimitiveSearchHTML(query, shroom.YieldHyphaNamesContainingString),
			u,
		),
	)
}
