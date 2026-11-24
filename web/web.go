// Package web contains web handlers and initialization stuff.
package web

import (
	"errors"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"net/url"

	"github.com/bouncepaw/mycorrhiza/help"
	"github.com/bouncepaw/mycorrhiza/history/histweb"
	"github.com/bouncepaw/mycorrhiza/hypview"
	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	"github.com/bouncepaw/mycorrhiza/internal/user"
	userdomain "github.com/bouncepaw/mycorrhiza/internal/user/domain"
	"github.com/bouncepaw/mycorrhiza/interwiki"
	"github.com/bouncepaw/mycorrhiza/misc"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/web/viewdata"

	"github.com/gorilla/mux"
)

// Handler initializes and returns the HTTP router based on the configuration.
func Handler() http.Handler {
	router := mux.NewRouter()
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
			util.PrepareRq(rq)
			w.Header().Add("Content-Security-Policy",
				"default-src 'self'; "+
					"img-src * data:; media-src *; style-src *; font-src * data:")
			next.ServeHTTP(w, rq)
		})
	})
	router.StrictSlash(true)

	// Public routes. They're always accessible regardless of the user status.
	misc.InitAssetHandlers(router)

	// Auth
	router.HandleFunc("/user-list", handlerUserList)
	router.HandleFunc("/lock", handlerLock)
	// The check below saves a lot of extra checks and lines of codes in other places in this file.
	if cfg.UseAuth {
		if cfg.AllowRegistration {
			router.HandleFunc("/register", handlerRegister).Methods(http.MethodPost, http.MethodGet)
		}
		router.HandleFunc("/login", handlerLogin)
		router.HandleFunc("/logout", handlerLogout)
	}

	// Wiki routes. They may be locked or restricted.
	r := router.PathPrefix("").Subrouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
			user := usersvc.FromRequest(rq)
			if !user.ShowLockMaybe(w, rq) {
				next.ServeHTTP(w, rq)
			}
		})
	})

	initReaders(r)
	initMutators(r)
	help.InitHandlers(r)
	misc.InitHandlers(r)
	hypview.Init()
	histweb.InitHandlers(r)
	interwikisvc.InitHandlers(r)

	r.PathPrefix("/add-to-category").HandlerFunc(handlerAddToCategory).Methods("POST")
	r.PathPrefix("/remove-from-category").HandlerFunc(handlerRemoveFromCategory).Methods("POST")
	r.PathPrefix("/category/").HandlerFunc(handlerCategory).Methods("GET")
	r.PathPrefix("/edit-category/").HandlerFunc(handlerEditCategory).Methods("GET")
	r.PathPrefix("/category").HandlerFunc(handlerListCategory).Methods("GET")

	// Admin routes
	if cfg.UseAuth {
		adminRouter := r.PathPrefix("/admin").Subrouter()
		adminRouter.Use(actionMiddleware(userdomain.ActionAdmin))

		adminRouter.HandleFunc("/shutdown", handlerAdminShutdown).Methods(http.MethodPost)
		adminRouter.HandleFunc("/reindex-users", handlerAdminReindexUsers).Methods(http.MethodPost)

		adminRouter.HandleFunc("/new-user", handlerAdminUserNew).Methods(http.MethodGet, http.MethodPost)
		adminRouter.HandleFunc("/users/{username}/edit", handlerAdminUserEdit).Methods(http.MethodGet, http.MethodPost)
		adminRouter.HandleFunc("/users/{username}/change-password", handlerAdminUserChangePassword).Methods(http.MethodPost)
		adminRouter.HandleFunc("/users/{username}/delete", handlerAdminUserDelete).Methods(http.MethodGet, http.MethodPost)
		adminRouter.HandleFunc("/users", handlerAdminUsers)

		adminRouter.HandleFunc("/", handlerAdmin)

		settingsRouter := r.PathPrefix("/settings").Subrouter()
		// TODO: check if necessary?
		//settingsRouter.Use(groupMiddleware("settings"))
		settingsRouter.HandleFunc("/change-password", handlerUserChangePassword).Methods(http.MethodGet, http.MethodPost)
	}

	// Index page
	r.HandleFunc("/", func(w http.ResponseWriter, rq *http.Request) {
		// Let's pray it never fails
		addr, _ := url.Parse("/hypha/" + cfg.HomeHypha)
		rq.URL = addr
		handlerHypha(w, rq)
	})

	initPages()

	return router
}

func actionMiddleware(action userdomain.Action) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
			if cfg.UseAuth && usersvc.CanProceed(rq, action) {
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

// Auth
func handlerUserList(w http.ResponseWriter, rq *http.Request) {
	admins, moderators, editors, readers := usersvc.UsersInGroups()
	_ = pageUserList.RenderTo(viewdata.MetaFrom(w, rq),
		map[string]any{
			"Admins":     admins,
			"Moderators": moderators,
			"Editors":    editors,
			"Readers":    readers,
		})
}

func handlerLock(w http.ResponseWriter, rq *http.Request) {
	_ = pageAuthLock.RenderTo(viewdata.MetaFrom(w, rq), map[string]any{})
}

// handlerRegister displays the register form (GET) or registers the user (POST).
func handlerRegister(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	if rq.Method == http.MethodGet {
		slog.Info("Showing registration form")
		_ = pageAuthRegister.RenderTo(viewdata.MetaFrom(w, rq), map[string]any{
			"UseAuth":           cfg.UseAuth,
			"AllowRegistration": cfg.AllowRegistration,
			"RawQuery":          rq.URL.RawQuery,
			"WikiName":          cfg.WikiName,
		})
		return
	}

	var (
		username = rq.PostFormValue("username")
		password = rq.PostFormValue("password")
		err      = usersvc.Register(username, password, userdomain.GroupEditor, userdomain.SourceLocal, false)
	)
	if err != nil {
		slog.Info("Failed to register", "username", username, "err", err.Error())
		w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
		w.WriteHeader(http.StatusBadRequest)
		_ = pageAuthRegister.RenderTo(viewdata.MetaFrom(w, rq), map[string]any{
			"UseAuth":           cfg.UseAuth,
			"AllowRegistration": cfg.AllowRegistration,
			"RawQuery":          rq.URL.RawQuery,
			"WikiName":          cfg.WikiName,

			"Err":      err,
			"Username": username,
			"Password": password,
		})
		return
	}

	slog.Info("Registered user", "username", username)
	if err := usersvc.LoginDataHTTP(w, username, password); err != nil {
		return
	}
	http.Redirect(w, rq, "/"+rq.URL.RawQuery, http.StatusSeeOther)
}

// handlerLogout shows the logout form (GET) or logs the user out (POST).
func handlerLogout(w http.ResponseWriter, rq *http.Request) {
	if rq.Method == http.MethodPost {
		slog.Info("Somebody logged out")
		usersvc.LogoutFromRequest(w, rq)
		http.Redirect(w, rq, "/", http.StatusSeeOther)
		return
	}

	var (
		u   = usersvc.FromRequest(rq)
		can = u != nil
	)
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	if can {
		slog.Info("Logging out", "username", u.Name)
		w.WriteHeader(http.StatusOK)
	} else {
		slog.Info("Unknown user logging out")
		w.WriteHeader(http.StatusForbidden)
	}
	_ = pageAuthLogout.RenderTo(viewdata.MetaFrom(w, rq), map[string]any{
		"CanLogout": can,
	})
}

// handlerLogin shows the login form (GET) or logs the user in (POST).
func handlerLogin(w http.ResponseWriter, rq *http.Request) {
	if rq.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		_ = pageAuthLogin.RenderTo(viewdata.MetaFrom(w, rq), map[string]any{
			"UseAuth":            cfg.UseAuth,
			"ErrUnknownUsername": false,
			"ErrWrongPassword":   false,
			"Err":                nil,
			"WikiName":           cfg.WikiName,
		})
		slog.Info("Somebody logging in")
		return
	}

	var (
		username = util.CanonicalName(rq.PostFormValue("username"))
		password = rq.PostFormValue("password")
		err      = usersvc.LoginDataHTTP(w, username, password)
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = pageAuthLogin.RenderTo(viewdata.MetaFrom(w, rq), map[string]any{
			"UseAuth":            cfg.UseAuth,
			"ErrUnknownUsername": errors.Is(err, usersvc.ErrUnknownUsername),
			"ErrWrongPassword":   errors.Is(err, usersvc.ErrWrongPassword),
			"Err":                err.Error(),
			"WikiName":           cfg.WikiName,
			"Username":           username,
		})
		slog.Info("Failed to log in", "username", username, "err", err.Error())
		return
	}
	http.Redirect(w, rq, "/", http.StatusSeeOther)
	slog.Info("Logged in", "username", username)
}
