package categories

import (
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strings"
)

// InitHandlers initializes HTTP handlers for the given router. Call somewhere in package web.
func InitHandlers(r *mux.Router) {
	r.PathPrefix("/add-to-category").HandlerFunc(handlerAddToCategory).Methods("POST")
	r.PathPrefix("/remove-from-category").HandlerFunc(handlerRemoveFromCategory).Methods("POST")
	r.PathPrefix("/category/").HandlerFunc(handlerCategory).Methods("GET")
	r.PathPrefix("/category").HandlerFunc(handlerListCategory).Methods("GET")
	prepareViews()
}

func handlerListCategory(w http.ResponseWriter, rq *http.Request) {
	log.Println("Viewing list of categories")
	categoryList(viewutil.MetaFrom(w, rq))
}

func handlerCategory(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	catName := util.CanonicalName(strings.TrimPrefix(strings.TrimPrefix(rq.URL.Path, "/category"), "/"))
	if catName == "" {
		handlerListCategory(w, rq)
		return
	}
	log.Println("Viewing category", catName)
	categoryPage(viewutil.MetaFrom(w, rq), catName)
}

func handlerRemoveFromCategory(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		hyphaName  = util.CanonicalName(rq.PostFormValue("hypha"))
		catName    = util.CanonicalName(rq.PostFormValue("cat"))
		redirectTo = rq.PostFormValue("redirect-to")
	)
	if !user.FromRequest(rq).CanProceed("remove-from-category") {
		w.WriteHeader(http.StatusForbidden)
		_, _ = io.WriteString(w, "403 Forbidden")
		return
	}
	if hyphaName == "" || catName == "" {
		http.Redirect(w, rq, redirectTo, http.StatusSeeOther)
		return
	}
	log.Println(user.FromRequest(rq).Name, "removed", hyphaName, "from", catName)
	removeHyphaFromCategory(hyphaName, catName)
	http.Redirect(w, rq, redirectTo, http.StatusSeeOther)
}

func handlerAddToCategory(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		hyphaName  = util.CanonicalName(rq.PostFormValue("hypha"))
		catName    = util.CanonicalName(rq.PostFormValue("cat"))
		redirectTo = rq.PostFormValue("redirect-to")
	)
	if !user.FromRequest(rq).CanProceed("add-to-category") {
		w.WriteHeader(http.StatusForbidden)
		_, _ = io.WriteString(w, "403 Forbidden")
		return
	}
	if hyphaName == "" || catName == "" {
		http.Redirect(w, rq, redirectTo, http.StatusSeeOther)
		return
	}
	log.Println(user.FromRequest(rq).Name, "added", hyphaName, "to", catName)
	addHyphaToCategory(hyphaName, catName)
	http.Redirect(w, rq, redirectTo, http.StatusSeeOther)
}
