package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

// GetRevision finds revision with id `id` of `hyphaName` in `hyphae`.
// If `id` is `"0"`, it means the last revision.
// If no such revision is found, last return value is false.
func GetRevision(hyphaName string, id string) (Revision, bool) {
	log.Println("Getting hypha", hyphaName, id)
	if hypha, ok := hyphae[hyphaName]; ok {
		if id == "0" {
			id = hypha.NewestRevision()
		}
		if rev, ok := hypha.Revisions[id]; ok {
			return *rev, true
		}
	}
	return Revision{}, false
}

// RevInMap finds value of `rev` (the one from URL queries like) in the passed map that is usually got from `mux.Vars(*http.Request)`.
// If there is no `rev`, return "0".
func RevInMap(m map[string]string) string {
	if id, ok := m["rev"]; ok {
		return id
	}
	return "0"
}

// `rootWikiDir` is a directory where all wiki files reside.
var rootWikiDir string

// `hyphae` is a map with all hyphae. Many functions use it.
var hyphae map[string]*Hypha

func main() {
	if len(os.Args) == 1 {
		panic("Expected a root wiki pages directory")
	}
	// Required so the rootWikiDir hereinbefore does not get redefined.
	var err error
	rootWikiDir, err = filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}

	log.Println("Welcome to MycorrhizaWiki α")
	log.Println("Indexing hyphae...")
	hyphae = recurFindHyphae(rootWikiDir)
	log.Println("Indexed", len(hyphae), "hyphae. Ready to accept requests.")

	// Start server code. See handlers.go for handlers' implementations.
	r := mux.NewRouter()

	r.Queries("action", "getBinary", "rev", revQuery).Path(hyphaUrl).
		HandlerFunc(HandlerGetBinary)
	r.Queries("action", "getBinary").Path(hyphaUrl).
		HandlerFunc(HandlerGetBinary)

	r.Queries("action", "raw", "rev", revQuery).Path(hyphaUrl).
		HandlerFunc(HandlerRaw)
	r.Queries("action", "raw").Path(hyphaUrl).
		HandlerFunc(HandlerRaw)

	r.Queries("action", "zen", "rev", revQuery).Path(hyphaUrl).
		HandlerFunc(HandlerZen)
	r.Queries("action", "zen").Path(hyphaUrl).
		HandlerFunc(HandlerZen)

	r.Queries("action", "view", "rev", revQuery).Path(hyphaUrl).
		HandlerFunc(HandlerView)
	r.Queries("action", "view").Path(hyphaUrl).
		HandlerFunc(HandlerView)

	r.Queries("action", "history").Path(hyphaUrl).
		HandlerFunc(HandlerHistory)

	r.Queries("action", "edit").Path(hyphaUrl).
		HandlerFunc(HandlerEdit)

	r.Queries("action", "rewind", "rev", revQuery).Path(hyphaUrl).
		HandlerFunc(HandlerRewind)

	r.Queries("action", "delete").Path(hyphaUrl).
		HandlerFunc(HandlerDelete)

	r.Queries("action", "rename", "to", hyphaPattern).Path(hyphaUrl).
		HandlerFunc(HandlerRename)

	r.Queries("action", "update").Path(hyphaUrl).Methods("POST").
		HandlerFunc(HandlerUpdate)

	r.HandleFunc(hyphaUrl, HandlerView)

	// Debug page that renders all hyphae.
	// TODO: make it redirect to home page.
	// TODO: make a home page.
	r.HandleFunc("/", func(w http.ResponseWriter, rq *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		for _, h := range hyphae {
			log.Println("Rendering latest revision of hypha", h.FullName)
			html, err := h.AsHtml("0", w)
			if err != nil {
				fmt.Fprintln(w, err)
			}
			fmt.Fprintln(w, html)
		}
	})

	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
