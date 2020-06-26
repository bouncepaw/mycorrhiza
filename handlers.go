package main

import (
	"log"
	// "io/ioutil"
	// "log"
	"net/http"
	// "path/filepath"
	// "strconv"
	// "strings"
	// "time"

	"github.com/bouncepaw/mycorrhiza/fs"
	"github.com/bouncepaw/mycorrhiza/render"
	"github.com/gorilla/mux"
)

// There are handlers below. See main() for their usage.

// Boilerplate code present in many handlers. Good to have it.
func HandlerBase(w http.ResponseWriter, rq *http.Request) (*fs.Hypha, bool) {
	vars := mux.Vars(rq)
	h := fs.Hs.Open(vars["hypha"]).OnRevision(RevInMap(vars))
	if h.Invalid {
		log.Println(h.Err)
		return h, false
	}
	return h, true
}

func HandlerRaw(w http.ResponseWriter, rq *http.Request) {
	log.Println("?action=raw")
	if h, ok := HandlerBase(w, rq); ok {
		h.ActionRaw(w)
	}
}

func HandlerBinary(w http.ResponseWriter, rq *http.Request) {
	log.Println("?action=binary")
	if h, ok := HandlerBase(w, rq); ok {
		h.ActionBinary(w)
	}
}

func HandlerZen(w http.ResponseWriter, rq *http.Request) {
	if h, ok := HandlerBase(w, rq); ok {
		h.ActionZen(w)
	}
}

func HandlerView(w http.ResponseWriter, rq *http.Request) {
	if h, ok := HandlerBase(w, rq); ok {
		h.ActionView(w, render.HyphaPage, render.Hypha404)
	}
}

func HandlerEdit(w http.ResponseWriter, rq *http.Request) {
	vars := mux.Vars(rq)
	h := fs.Hs.Open(vars["hypha"]).OnRevision("0")
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(render.HyphaEdit(h)))
}

func HandlerUpdate(w http.ResponseWriter, rq *http.Request) {
	vars := mux.Vars(rq)
	log.Println("Attempt to update hypha", vars["hypha"])
	h := fs.Hs.
		Open(vars["hypha"]).
		CreateDirIfNeeded().
		AddRevisionFromHttpData(rq).
		WriteTextFileFromHttpData(rq).
		WriteBinaryFileFromHttpData(rq).
		SaveJson().
		Store()

	if h.Invalid {
		log.Println(h.Err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(render.HyphaUpdateOk(h)))
}
