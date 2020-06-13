package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	// "strconv"
	"time"
)

var rootWikiDir string

func hyphaeAsMap(hyphae []*Hypha) map[string]*Hypha {
	mh := make(map[string]*Hypha)
	for _, h := range hyphae {
		mh[h.Name] = h
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

	hyphae := hyphaeAsMap(recurFindHyphae(rootWikiDir))
	setRelations(hyphae)

	// Start server code
	r := mux.NewRouter()
	r.HandleFunc("/showHyphae", func(w http.ResponseWriter, r *http.Request) {
		for _, h := range hyphae {
			fmt.Fprintln(w, h)
		}
	})
	r.Queries(
		"action", "getBinary",
		"rev", "{rev:[\\d]+}",
	).Path("/{hypha:" + hyphaPattern + "}").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			rev, err := GetRevision(hyphae, vars["hypha"], vars["rev"])
			if err != nil {
				log.Println("Failed to show image of", rev.FullName)
			}
			fileContents, err := ioutil.ReadFile(rev.BinaryPath)
			if err != nil {
				log.Println("Failed to show image of", rev.FullName)
			}
			log.Println("Contents:", fileContents[:10], "...")
			w.Header().Set("Content-Type", rev.MimeType)
			// w.Header().Set("Content-Length", strconv.Itoa(len(fileContents)))
			w.WriteHeader(http.StatusOK)
			w.Write(fileContents)
			log.Println("Showing image of", rev.FullName)
		})

	r.HandleFunc("/{hypha:"+hyphaPattern+"}",
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			rev, err := GetRevision(hyphae, vars["hypha"], "0")
			if err != nil {
				log.Println("Failed to show image of", rev.FullName)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			html, err := rev.Render(hyphae)
			if err != nil {
				log.Println("Failed to show image of", rev.FullName)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, HyphaPage(hyphae, rev, html))
			log.Println("Rendering", rev.FullName)
		})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		for _, v := range hyphae {
			log.Println("Rendering latest revision of hypha", v.Name)
			html, err := v.Render(hyphae, 0)
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
