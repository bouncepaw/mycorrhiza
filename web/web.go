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
	router := mux.NewRouter()
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Do stuff here
			log.Println(r.RequestURI)
			// Call the next handler, which can be another middleware in the chain, or the final handler.
			next.ServeHTTP(w, r)
		})
	})

	// Public routes
	initAuth(router)
	router.HandleFunc("/robots.txt", handlerRobotsTxt)
	router.HandleFunc("/static/style.css", handlerStyle)
	router.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))))

	// Wiki routes. They may be locked or restricted
	wikiRouter := router.PathPrefix("").Subrouter()
	wikiRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
			user := user.FromRequest(rq)
			if !user.ShowLockMaybe(w, rq) {
				next.ServeHTTP(w, rq)
			}
		})
	})

	initReaders(wikiRouter)
	initMutators(wikiRouter)
	initAdmin(wikiRouter)
	initHistory(wikiRouter)
	initStuff(wikiRouter)
	initSearch(wikiRouter)

	// Miscellaneous
	wikiRouter.HandleFunc("/user-list", handlerUserList)

	// Index page
	wikiRouter.HandleFunc("/", func(w http.ResponseWriter, rq *http.Request) {
		// Let's pray it never fails
		addr, _ := url.Parse("/hypha/" + cfg.HomeHypha)
		rq.URL = addr
		handlerHypha(w, rq)
	})

	return router
}
