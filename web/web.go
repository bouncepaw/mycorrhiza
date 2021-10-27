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
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/static"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"
)

var stylesheets = []string{"default.css", "custom.css"}

// httpErr is used by many handlers to signal errors in a compact way.
func httpErr(w http.ResponseWriter, lc *l18n.Localizer, status int, name, title, errMsg string) {
	log.Println(errMsg, "for", name)
	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	w.WriteHeader(status)
	_, err := fmt.Fprint(
		w,
		views.BaseHTML(
			title,
			fmt.Sprintf(
				`<main class="main-width"><p>%s. <a href="/hypha/%s">%s<a></p></main>`,
				errMsg,
				name,
				lc.Get("ui.error_go_back"),
			),
			lc,
			user.EmptyUser(),
		),
	)
	if err != nil {
		log.Println("an error occurred in httpErr function:", err)
	}
}

func handlerStyle(w http.ResponseWriter, rq *http.Request) {
	w.Header().Set("Content-Type", mime.TypeByExtension(".css"))
	for _, name := range stylesheets {
		file, err := static.FS.Open(name)
		if err != nil {
			log.Println("an error occurred in handlerStyle function:", err)
			continue
		}
		_, err = io.Copy(w, file)
		if err != nil {
			log.Println("an error occurred in handlerStyle function:", err)
			continue
		}
		err = file.Close()
		if err != nil {
			log.Println("an error occurred in handlerStyle function:", err)
			continue
		}
	}
}

func handlerUserList(w http.ResponseWriter, rq *http.Request) {
	lc := l18n.FromRequest(rq)
	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(views.BaseHTML(lc.Get("ui.users_title"), views.UserListHTML(lc), lc, user.FromRequest(rq))))
	if err != nil {
		log.Println("an error occurred in handlerUserList function:", err)
	}
}

func handlerRobotsTxt(w http.ResponseWriter, rq *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	file, err := static.FS.Open("robots.txt")
	if err != nil {
		log.Println("an error occurred in handlerRobotsTxt function:", err)
		return
	}
	_, err = io.Copy(w, file)
	// Even if we failed copying stuff into the response writer we can try close file anyway
	if err != nil {
		log.Println("an error occurred in handlerRobotsTxt function:", err)
	}
	err = file.Close()
	if err != nil {
		log.Println("an error occurred in handlerRobotsTxt function:", err)
		return
	}
}

// Handler initializes and returns the HTTP router based on the configuration.
func Handler() http.Handler {
	router := mux.NewRouter()
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
			util.PrepareRq(rq)
			next.ServeHTTP(w, rq)
		})
	})
	router.StrictSlash(true)

	// Public routes. They're always accessible regardless of the user status.
	initAuth(router)
	router.HandleFunc("/robots.txt", handlerRobotsTxt)
	router.HandleFunc("/static/style.css", handlerStyle)
	router.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))))

	// Wiki routes. They may be locked or restricted.
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
	initHistory(wikiRouter)
	initStuff(wikiRouter)
	initSearch(wikiRouter)
	initBacklinks(wikiRouter)

	// Admin routes.
	if cfg.UseAuth {
		adminRouter := wikiRouter.PathPrefix("/admin").Subrouter()
		adminRouter.Use(groupMiddleware("admin"))
		initAdmin(adminRouter)
	}

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

func groupMiddleware(group string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
			if cfg.UseAuth && user.CanProceed(rq, group) {
				next.ServeHTTP(w, rq)
				return
			}

			// TODO: handle this better. Merge this code with all other
			// authorization code in this project.

			w.WriteHeader(http.StatusForbidden)
			_, err := io.WriteString(w, "403 forbidden")
			if err != nil {
				log.Println("an error occurred in groupMiddleware function:", err)
			}
		})
	}
}
