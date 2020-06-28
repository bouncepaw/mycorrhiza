package main

import (
	"log"
	"net/http"

	"github.com/bouncepaw/mycorrhiza/fs"
	"github.com/bouncepaw/mycorrhiza/render"
	"github.com/gorilla/mux"
)

// There are handlers below. See main() for their usage.

// Boilerplate code present in many handlers. Good to have it.
func HandlerBase(w http.ResponseWriter, rq *http.Request) *fs.Hypha {
	vars := mux.Vars(rq)
	return fs.Hs.Open(vars["hypha"]).OnRevision(RevInMap(vars))
}

func HandlerRaw(w http.ResponseWriter, rq *http.Request) {
	log.Println("?action=raw")
	HandlerBase(w, rq).ActionRaw(w).LogSuccMaybe("Serving raw text")
}

func HandlerBinary(w http.ResponseWriter, rq *http.Request) {
	log.Println("?action=binary")
	HandlerBase(w, rq).ActionBinary(w).LogSuccMaybe("Serving binary data")
}

func HandlerZen(w http.ResponseWriter, rq *http.Request) {
	log.Println("?action=zen")
	HandlerBase(w, rq).ActionZen(w).LogSuccMaybe("Rendering zen")
}

func HandlerView(w http.ResponseWriter, rq *http.Request) {
	log.Println("?action=view")
	HandlerBase(w, rq).
		ActionView(w, render.HyphaPage, render.Hypha404).
		LogSuccMaybe("Rendering hypha view")
}

func HandlerEdit(w http.ResponseWriter, rq *http.Request) {
	vars := mux.Vars(rq)
	h := fs.Hs.Open(vars["hypha"]).OnRevision("0")
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(render.HyphaEdit(h))
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
		Store().
		LogSuccMaybe("Saved changes")

	if !h.Invalid {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(render.HyphaUpdateOk(h))
	}
}
