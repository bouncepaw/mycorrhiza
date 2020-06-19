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

func GetRevision(hyphae map[string]*Hypha, hyphaName string, rev string, w http.ResponseWriter) (Revision, bool) {
	log.Println("Getting hypha", hyphaName, rev)
	for name, hypha := range hyphae {
		if name == hyphaName {
			if rev == "0" {
				rev = hypha.NewestRevision()
			}
			for id, r := range hypha.Revisions {
				if rev == id {
					return *r, true
				}
			}
		}
	}
	return Revision{}, false
}

func RevInMap(m map[string]string) string {
	if val, ok := m["rev"]; ok {
		return val
	}
	return "0"
}

var rootWikiDir string
var hyphae map[string]*Hypha

func hyphaeAsMap(hyphae []*Hypha) map[string]*Hypha {
	mh := make(map[string]*Hypha)
	for _, h := range hyphae {
		mh[h.Name()] = h
	}
	return mh
}

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
	// setRelations(hyphae)

	// Start server code
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

	r.Queries(
		"action", "update",
	).Path(hyphaUrl).Methods("POST").
		HandlerFunc(HandlerUpdate)

	r.HandleFunc(hyphaUrl, HandlerView)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		for _, v := range hyphae {
			log.Println("Rendering latest revision of hypha", v.Name())
			html, err := v.AsHtml(hyphae, "0")
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
