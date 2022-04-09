// Package web contains web handlers and initialization stuff.
package web

import (
	"github.com/bouncepaw/mycorrhiza/backlinks"
	"github.com/bouncepaw/mycorrhiza/categories"
	"github.com/bouncepaw/mycorrhiza/help"
	"github.com/bouncepaw/mycorrhiza/misc"
	"io"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
)

// Handler initializes and returns the HTTP router based on the configuration.
func Handler() http.Handler {
	router := mux.NewRouter()
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
			util.PrepareRq(rq)
			w.Header().Add("Content-Security-Policy",
				"default-src 'self' telegram.org *.telegram.org; "+
					"img-src * data:; media-src *; style-src *; font-src * data:")
			next.ServeHTTP(w, rq)
		})
	})
	router.StrictSlash(true)

	// Public routes. They're always accessible regardless of the user status.
	initAuth(router)

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
	help.InitHandlers(wikiRouter)
	backlinks.InitHandlers(wikiRouter)
	categories.InitHandlers(wikiRouter)
	misc.InitHandlers(wikiRouter)

	// Admin routes.
	if cfg.UseAuth {
		adminRouter := wikiRouter.PathPrefix("/admin").Subrouter()
		adminRouter.Use(groupMiddleware("admin"))
		initAdmin(adminRouter)
	}

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
			io.WriteString(w, "403 forbidden")
		})
	}
}
