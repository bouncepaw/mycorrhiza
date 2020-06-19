package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Boilerplate code present in many handlers. Good to have it.
func HandlerBase(w http.ResponseWriter, r *http.Request) (Revision, bool) {
	vars := mux.Vars(r)
	revno := RevInMap(vars)
	return GetRevision(hyphae, vars["hypha"], revno, w)
}

func HandlerGetBinary(w http.ResponseWriter, r *http.Request) {
	if rev, ok := HandlerBase(w, r); ok {
		rev.ActionGetBinary(w)
	}
}

func HandlerRaw(w http.ResponseWriter, r *http.Request) {
	if rev, ok := HandlerBase(w, r); ok {
		rev.ActionRaw(w)
	}
}

func HandlerZen(w http.ResponseWriter, r *http.Request) {
	if rev, ok := HandlerBase(w, r); ok {
		rev.ActionZen(w)
	}
}

func HandlerView(w http.ResponseWriter, r *http.Request) {
	if rev, ok := HandlerBase(w, r); ok {
		rev.ActionView(w, HyphaPage)
	}
}

func HandlerHistory(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	log.Println("Attempt to access an unimplemented thing")
}

func HandlerEdit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ActionEdit(vars["hypha"], w)
}

func HandlerRewind(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	log.Println("Attempt to access an unimplemented thing")
}

func HandlerDelete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	log.Println("Attempt to access an unimplemented thing")
}

func HandlerRename(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	log.Println("Attempt to access an unimplemented thing")
}

func HandlerUpdate(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	log.Println("Attempt to access an unimplemented thing")
}
