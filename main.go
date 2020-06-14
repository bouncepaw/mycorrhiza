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

func RevInMap(m map[string]string) string {
	if val, ok := m["rev"]; ok {
		return val
	}
	return "0"
}

// handlers
func HandlerGetBinary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	revno := RevInMap(vars)
	rev, ok := GetRevision(hyphae, vars["hypha"], revno, w)
	if !ok {
		return
	}
	fileContents, err := ioutil.ReadFile(rev.BinaryPath)
	if err != nil {
		log.Println("Failed to load binary data of", rev.FullName, rev.Id)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", rev.MimeType)
	w.WriteHeader(http.StatusOK)
	w.Write(fileContents)
	log.Println("Showing image of", rev.FullName, rev.Id)
}

func HandlerRaw(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	revno := RevInMap(vars)
	rev, ok := GetRevision(hyphae, vars["hypha"], revno, w)
	if !ok {
		return
	}
	fileContents, err := ioutil.ReadFile(rev.TextPath)
	if err != nil {
		log.Println("Failed to load text data of", rev.FullName, rev.Id)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(fileContents)
	log.Println("Serving text data of", rev.FullName, rev.Id)
}

func HandlerZen(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	revno := RevInMap(vars)
	rev, ok := GetRevision(hyphae, vars["hypha"], revno, w)
	if !ok {
		return
	}
	html, err := rev.Render(hyphae)
	if err != nil {
		log.Println("Failed to render", rev.FullName)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, html)
}

func HandlerView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	revno := RevInMap(vars)
	rev, ok := GetRevision(hyphae, vars["hypha"], revno, w)
	if !ok {
		return
	}
	html, err := rev.Render(hyphae)
	if err != nil {
		log.Println("Failed to render", rev.FullName)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, HyphaPage(hyphae, rev, html))
	log.Println("Rendering", rev.FullName)
}

func HandlerHistory(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	log.Println("Attempt to access an unimplemented thing")
}

func HandlerEdit(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	log.Println("Attempt to access an unimplemented thing")
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

var rootWikiDir string
var hyphae map[string]*Hypha

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

	hyphae = hyphaeAsMap(recurFindHyphae(rootWikiDir))
	setRelations(hyphae)

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

	r.Queries("action", "update").Path(hyphaUrl).
		HandlerFunc(HandlerUpdate)

	r.HandleFunc(hyphaUrl, HandlerView)

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
