package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/fs"
	"github.com/gorilla/mux"
)

// RevInMap finds value of `rev` (the one from URL queries like) in the passed map that is usually got from `mux.Vars(*http.Request)`.
// If there is no `rev`, return "0".
func RevInMap(m map[string]string) string {
	if id, ok := m["rev"]; ok {
		return id
	}
	return "0"
}

func IdempotentRouterBoiler(router *mux.Router, action string, handler func(w http.ResponseWriter, rq *http.Request)) {
	router.
		Queries("action", action, "rev", cfg.RevQuery).
		Path(cfg.MyceliumUrl + cfg.HyphaUrl).
		HandlerFunc(handler)
	router.
		Queries("action", action).
		Path(cfg.MyceliumUrl + cfg.HyphaUrl).
		HandlerFunc(handler)
	router.
		Queries("action", action, "rev", cfg.RevQuery).
		Path(cfg.HyphaUrl).
		HandlerFunc(handler)
	router.
		Queries("action", action).
		Path(cfg.HyphaUrl).
		HandlerFunc(handler)
}

func main() {
	if len(os.Args) == 1 {
		panic("Expected a root wiki pages directory")
	}
	wikiDir, err := filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}

	log.Println("Welcome to MycorrhizaWiki α")
	cfg.InitConfig(wikiDir)
	log.Println("Indexing hyphae...")
	fs.VerifyMycelia()
	fs.InitStorage()

	// Start server code. See handlers.go for handlers' implementations.
	r := mux.NewRouter()

	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, rq *http.Request) {
		http.ServeFile(w, rq, filepath.Join(filepath.Dir(cfg.WikiDir), "favicon.ico"))
	})

	IdempotentRouterBoiler(r, "binary", HandlerBinary)
	IdempotentRouterBoiler(r, "raw", HandlerRaw)
	IdempotentRouterBoiler(r, "zen", HandlerZen)
	IdempotentRouterBoiler(r, "view", HandlerView)

	r.Queries("action", "edit").Path(cfg.MyceliumUrl + cfg.HyphaUrl).
		HandlerFunc(HandlerEdit)
	r.Queries("action", "edit").Path(cfg.HyphaUrl).
		HandlerFunc(HandlerEdit)

	r.Queries("action", "update").Path(cfg.MyceliumUrl + cfg.HyphaUrl).
		Methods("POST").HandlerFunc(HandlerUpdate)
	r.Queries("action", "update").Path(cfg.HyphaUrl).
		Methods("POST").HandlerFunc(HandlerUpdate)

	r.HandleFunc(cfg.HyphaUrl, HandlerView)

	// Debug page that renders all hyphae.
	// TODO: make it redirect to home page.
	// TODO: make a home page.
	r.HandleFunc("/", func(w http.ResponseWriter, rq *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `
<p>Check out <a href="/Fruit">Fruit</a> maybe.</p>
<p><strong>TODO:</strong> make this page usable</p>
		`)
	})

	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    cfg.Address,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
