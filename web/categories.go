package web

import (
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"
	"github.com/gorilla/mux"
	"net/http"
)

func initCategories(r *mux.Router) {
	r.PathPrefix("/add-to-category").HandlerFunc(handlerAddToCategory).Methods("POST")
	r.PathPrefix("/remove-from-category").HandlerFunc(handlerRemoveFromCategory).Methods("POST")
	r.PathPrefix("/category/").HandlerFunc(handlerCategory).Methods("GET")
}

func handlerCategory(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		catName = util.HyphaNameFromRq(rq, "category")
	)
	views.CategoryPageHTML(views.MetaFrom(w, rq), catName)
}

func handlerRemoveFromCategory(w http.ResponseWriter, rq *http.Request) {

}

func handlerAddToCategory(w http.ResponseWriter, rq *http.Request) {

}
