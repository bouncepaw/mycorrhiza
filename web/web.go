// Package web contains web handlers and initialization stuff.
//
// It exports just one function: Init. Call it if you want to have web capabilities.
package web

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/static"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"
)

var stylesheets = []string{"default.css", "custom.css"}

// httpErr is used by many handlers to signal errors in a compact way.
func httpErr(w http.ResponseWriter, status int, name, title, errMsg string) {
	log.Println(errMsg, "for", name)
	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	w.WriteHeader(status)
	fmt.Fprint(
		w,
		views.BaseHTML(
			title,
			fmt.Sprintf(
				`<main class="main-width"><p>%s. <a href="/hypha/%s">Go back to the hypha.<a></p></main>`,
				errMsg,
				name,
			),
			user.EmptyUser(),
		),
	)
}

func handlerStyle(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)

	w.Header().Set("Content-Type", mime.TypeByExtension(".css"))
	for _, name := range stylesheets {
		file, err := static.FS.Open(name)
		if err != nil {
			continue
		}
		io.Copy(w, file)
		file.Close()
	}
}

func handlerUserList(w http.ResponseWriter, rq *http.Request) {
	if shown := user.FromRequest(rq).ShowLockMaybe(w, rq); shown {
		return
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(views.BaseHTML("User list", views.UserListHTML(), user.FromRequest(rq))))
}

func handlerRobotsTxt(w http.ResponseWriter, rq *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	file, err := static.FS.Open("robots.txt")
	if err != nil {
		return
	}
	io.Copy(w, file)
	file.Close()
}

func Handler() http.Handler {
	r := mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Do stuff here
			log.Println(r.RequestURI)
			// Call the next handler, which can be another middleware in the chain, or the final handler.
			next.ServeHTTP(w, r)
		})
	})

	// Available all the time
	initAuth(r)

	initReaders(r)
	initMutators(r)

	initAdmin(r)
	initHistory(r)
	initStuff(r)
	initSearch(r)

	// Miscellaneous
	r.HandleFunc("/user-list", handlerUserList)
	r.HandleFunc("/robots.txt", handlerRobotsTxt)

	// Static assets
	r.HandleFunc("/static/style.css", handlerStyle)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))))

	// Index page
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Let's pray it never fails
		addr, _ := url.Parse("/hypha/" + cfg.HomeHypha)
		r.URL = addr
		handlerHypha(w, r)
	})

	return r
}
