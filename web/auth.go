package web

import (
	"errors"
	"fmt"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"io"
	"log"
	"mime"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"
)

func initAuth(r *mux.Router) {
	r.HandleFunc("/user-list", handlerUserList)
	r.HandleFunc("/lock", handlerLock)
	// The check below saves a lot of extra checks and lines of codes in other places in this file.
	if !cfg.UseAuth {
		return
	}
	if cfg.AllowRegistration {
		r.HandleFunc("/register", handlerRegister)
	}
	if cfg.TelegramEnabled {
		r.HandleFunc("/telegram-login", handlerTelegramLogin)
	}
	r.HandleFunc("/login", handlerLogin)
	r.HandleFunc("/logout", handlerLogout)
}

func handlerUserList(w http.ResponseWriter, rq *http.Request) {
	lc := l18n.FromRequest(rq)
	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(views.Base(viewutil.MetaFrom(w, rq), lc.Get("ui.users_title"), views.UserList(lc))))
}

func handlerLock(w http.ResponseWriter, rq *http.Request) {
	_, _ = io.WriteString(w, views.Lock(l18n.FromRequest(rq)))
}

// handlerRegister displays the register form (GET) or registers the user (POST).
func handlerRegister(w http.ResponseWriter, rq *http.Request) {
	lc := l18n.FromRequest(rq)
	util.PrepareRq(rq)
	if rq.Method == http.MethodGet {
		_, _ = io.WriteString(
			w,
			views.Base(
				viewutil.MetaFrom(w, rq),
				lc.Get("auth.register_title"),
				views.Register(rq),
			),
		)
	} else if rq.Method == http.MethodPost {
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
				views.Base(
					viewutil.MetaFrom(w, rq),
					lc.Get("auth.register_title"),
					fmt.Sprintf(
						`<main class="main-width"><p>%s</p><p><a href="/register">%s<a></p></main>`,
						err.Error(),
						lc.Get("auth.try_again"),
					),
				),
			)
		} else {
			log.Printf("Successfully registered ‘%s’", username)
			user.LoginDataHTTP(w, rq, username, password)
			http.Redirect(w, rq, "/"+rq.URL.RawQuery, http.StatusSeeOther)
		}
	}
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
			views.Base(viewutil.MetaFrom(w, rq), lc.Get("auth.logout_title"), views.Logout(can, lc)),
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
			views.Base(
				viewutil.MetaFrom(w, rq),
				lc.Get("auth.login_title"),
				views.Login(lc),
			),
		)
	} else if rq.Method == http.MethodPost {
		var (
			username = util.CanonicalName(rq.PostFormValue("username"))
			password = rq.PostFormValue("password")
			err      = user.LoginDataHTTP(w, rq, username, password)
		)
		if err != "" {
			w.Header().Set("Content-Type", "text/html;charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = io.WriteString(w, views.Base(viewutil.MetaFrom(w, rq), err, views.LoginError(err, lc)))
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
	if user.HasUsername(username) && user.ByName(username).Source == "telegram" {
		// Problems is something we put blankets on.
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
			views.Base(
				viewutil.MetaFrom(w, rq),
				lc.Get("ui.error"),
				fmt.Sprintf(
					`<main class="main-width"><p>%s</p><p>%s</p><p><a href="/login">%s<a></p></main>`,
					lc.Get("auth.error_telegram"),
					err.Error(),
					lc.Get("auth.go_login"),
				),
			),
		)
		return
	}

	errmsg := user.LoginDataHTTP(w, rq, username, "")
	if errmsg != "" {
		log.Printf("Failed to login ‘%s’ using Telegram: %s", username, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(
			w,
			views.Base(
				viewutil.MetaFrom(w, rq),
				"Error",
				fmt.Sprintf(
					`<main class="main-width"><p>%s</p><p>%s</p><p><a href="/login">%s<a></p></main>`,
					lc.Get("auth.error_telegram"),
					err.Error(),
					lc.Get("auth.go_login"),
				),
			),
		)
		return
	}
	log.Printf("Authorize ‘%s’ from Telegram", username)
	http.Redirect(w, rq, "/", http.StatusSeeOther)
}
