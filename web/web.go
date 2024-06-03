// Package web contains web handlers and initialization stuff.
package web

import (
	"errors"
	"fmt"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"strings"

	"github.com/bouncepaw/mycorrhiza/admin"
	"github.com/bouncepaw/mycorrhiza/auth"
	"github.com/bouncepaw/mycorrhiza/categories"
	"github.com/bouncepaw/mycorrhiza/help"
	"github.com/bouncepaw/mycorrhiza/history/histweb"
	"github.com/bouncepaw/mycorrhiza/hypview"
	"github.com/bouncepaw/mycorrhiza/interwiki"
	"github.com/bouncepaw/mycorrhiza/misc"
	"github.com/bouncepaw/mycorrhiza/settings"

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
	misc.InitAssetHandlers(router)

	// Auth
	router.HandleFunc("/user-list", handlerUserList)
	router.HandleFunc("/lock", handlerLock)
	// The check below saves a lot of extra checks and lines of codes in other places in this file.
	if cfg.UseAuth {
		if cfg.AllowRegistration {
			router.HandleFunc("/register", handlerRegister).Methods(http.MethodPost, http.MethodGet)
		}
		if cfg.TelegramEnabled {
			router.HandleFunc("/telegram-login", handlerTelegramLogin)
		}
		router.HandleFunc("/login", handlerLogin)
		router.HandleFunc("/logout", handlerLogout)
	}

	// Wiki routes. They may be locked or restricted.
	r := router.PathPrefix("").Subrouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
			user := user.FromRequest(rq)
			if !user.ShowLockMaybe(w, rq) {
				next.ServeHTTP(w, rq)
			}
		})
	})

	initReaders(r)
	initMutators(r)
	help.InitHandlers(r)
	categories.InitHandlers(r)
	misc.InitHandlers(r)
	hypview.Init()
	histweb.InitHandlers(r)
	interwiki.InitHandlers(r)

	// Admin routes
	if cfg.UseAuth {
		adminRouter := r.PathPrefix("/admin").Subrouter()
		adminRouter.Use(groupMiddleware("admin"))
		admin.Init(adminRouter)

		settingsRouter := r.PathPrefix("/settings").Subrouter()
		// TODO: check if necessary?
		//settingsRouter.Use(groupMiddleware("settings"))
		settings.Init(settingsRouter)
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

// Auth
func handlerUserList(w http.ResponseWriter, rq *http.Request) {
	lc := l18n.FromRequest(rq)
	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(viewutil.Base(viewutil.MetaFrom(w, rq), lc.Get("ui.users_title"), auth.UserList(lc), map[string]string{})))
}

func handlerLock(w http.ResponseWriter, rq *http.Request) {
	_, _ = io.WriteString(w, auth.Lock(l18n.FromRequest(rq)))
}

// handlerRegister displays the register form (GET) or registers the user (POST).
func handlerRegister(w http.ResponseWriter, rq *http.Request) {
	lc := l18n.FromRequest(rq)
	util.PrepareRq(rq)
	if rq.Method == http.MethodGet {
		_, _ = io.WriteString(
			w,
			viewutil.Base(
				viewutil.MetaFrom(w, rq),
				lc.Get("auth.register_title"),
				auth.Register(rq),
				map[string]string{},
			),
		)
		return
	}

	var (
		username = rq.PostFormValue("username")
		password = rq.PostFormValue("password")
		err      = user.Register(username, password, "editor", "local", false)
	)
	if err != nil {
		log.Printf("Failed to register ‘%s’: %s", username, err.Error())
		w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(
			w,
			viewutil.Base(
				viewutil.MetaFrom(w, rq),
				lc.Get("auth.register_title"),
				fmt.Sprintf(
					`<main class="main-width"><p>%s</p><p><a href="/register">%s<a></p></main>`,
					err.Error(),
					lc.Get("auth.try_again"),
				),
				map[string]string{},
			),
		)
		return
	}

	log.Printf("Successfully registered ‘%s’", username)
	if err := user.LoginDataHTTP(w, username, password); err != nil {
		return
	}
	http.Redirect(w, rq, "/"+rq.URL.RawQuery, http.StatusSeeOther)
}

// handlerLogout shows the logout form (GET) or logs the user out (POST).
func handlerLogout(w http.ResponseWriter, rq *http.Request) {
	if rq.Method == http.MethodGet {
		var (
			u   = user.FromRequest(rq)
			can = u != nil
			lc  = l18n.FromRequest(rq)
		)
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		if can {
			log.Println("User", u.Name, "tries to log out")
			w.WriteHeader(http.StatusOK)
		} else {
			log.Println("Unknown user tries to log out")
			w.WriteHeader(http.StatusForbidden)
		}
		_, _ = io.WriteString(
			w,
			viewutil.Base(viewutil.MetaFrom(w, rq), lc.Get("auth.logout_title"), auth.Logout(can, lc), map[string]string{}),
		)
	} else if rq.Method == http.MethodPost {
		user.LogoutFromRequest(w, rq)
		http.Redirect(w, rq, "/", http.StatusSeeOther)
	}
}

// handlerLogin shows the login form (GET) or logs the user in (POST).
func handlerLogin(w http.ResponseWriter, rq *http.Request) {
	lc := l18n.FromRequest(rq)
	if rq.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(
			w,
			viewutil.Base(
				viewutil.MetaFrom(w, rq),
				lc.Get("auth.login_title"),
				auth.Login(lc),
				map[string]string{},
			),
		)
	} else if rq.Method == http.MethodPost {
		var (
			username = util.CanonicalName(rq.PostFormValue("username"))
			password = rq.PostFormValue("password")
			err      = user.LoginDataHTTP(w, username, password)
		)
		if err != nil {
			w.Header().Set("Content-Type", "text/html;charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = io.WriteString(w, viewutil.Base(viewutil.MetaFrom(w, rq), err.Error(), auth.LoginError(err.Error(), lc), map[string]string{}))
			return
		}
		http.Redirect(w, rq, "/", http.StatusSeeOther)
	}
}

func handlerTelegramLogin(w http.ResponseWriter, rq *http.Request) {
	// Note there is no lock here.
	lc := l18n.FromRequest(rq)
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	rq.ParseForm()
	var (
		values     = rq.URL.Query()
		username   = strings.ToLower(values.Get("username"))
		seemsValid = user.TelegramAuthParamsAreValid(values)
		err        = user.Register(
			username,
			"", // Password matters not
			"editor",
			"telegram",
			false,
		)
	)
	// If registering a user via Telegram failed, because a Telegram user with this name
	// has already registered, then everything is actually ok!
	if user.HasUsername(username) && user.ByName(username).Source == "telegram" {
		err = nil
	}

	if !seemsValid {
		err = errors.New("Wrong parameters")
	}

	if err != nil {
		log.Printf("Failed to register ‘%s’ using Telegram: %s", username, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(
			w,
			viewutil.Base(
				viewutil.MetaFrom(w, rq),
				lc.Get("ui.error"),
				fmt.Sprintf(
					`<main class="main-width"><p>%s</p><p>%s</p><p><a href="/login">%s<a></p></main>`,
					lc.Get("auth.error_telegram"),
					err.Error(),
					lc.Get("auth.go_login"),
				),
				map[string]string{},
			),
		)
		return
	}

	errmsg := user.LoginDataHTTP(w, username, "")
	if errmsg != nil {
		log.Printf("Failed to login ‘%s’ using Telegram: %s", username, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(
			w,
			viewutil.Base(
				viewutil.MetaFrom(w, rq),
				"Error",
				fmt.Sprintf(
					`<main class="main-width"><p>%s</p><p>%s</p><p><a href="/login">%s<a></p></main>`,
					lc.Get("auth.error_telegram"),
					err.Error(),
					lc.Get("auth.go_login"),
				),
				map[string]string{},
			),
		)
		return
	}
	log.Printf("Authorize ‘%s’ from Telegram", username)
	http.Redirect(w, rq, "/", http.StatusSeeOther)
}
