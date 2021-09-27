package web

import (
	"errors"
	"fmt"
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
	r.HandleFunc("/lock", handlerLock)
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
	r.HandleFunc("/login-data", handlerLoginData)
	r.HandleFunc("/logout", handlerLogout)
	r.HandleFunc("/logout-confirm", handlerLogoutConfirm)
}

func handlerLock(w http.ResponseWriter, rq *http.Request) {
	io.WriteString(w, views.LockHTML(l18n.FromRequest(rq)))
}

// handlerRegister both displays the register form (GET) and registers users (POST).
func handlerRegister(w http.ResponseWriter, rq *http.Request) {
	lc := l18n.FromRequest(rq)
	util.PrepareRq(rq)
	if !cfg.AllowRegistration {
		w.WriteHeader(http.StatusForbidden)
	}
	if rq.Method == http.MethodGet {
		io.WriteString(
			w,
			views.BaseHTML(
				lc.Get("auth.register_title"),
				views.RegisterHTML(rq),
				lc,
				user.FromRequest(rq),
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
			fmt.Fprint(
				w,
				views.BaseHTML(
					lc.Get("auth.register_title"),
					fmt.Sprintf(
						`<main class="main-width"><p>%s</p><p><a href="/register">%s<a></p></main>`,
						err.Error(),
						lc.Get("auth.try_again"),
					),
					lc,
					user.FromRequest(rq),
				),
			)
		} else {
			log.Printf("Successfully registered ‘%s’", username)
			user.LoginDataHTTP(w, rq, username, password)
			http.Redirect(w, rq, "/"+rq.URL.RawQuery, http.StatusSeeOther)
		}
	}
}

// handlerLogout shows the logout form.
func handlerLogout(w http.ResponseWriter, rq *http.Request) {
	var (
		u   = user.FromRequest(rq)
		can = u != nil
		lc = l18n.FromRequest(rq)
	)
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	if can {
		log.Println("User", u.Name, "tries to log out")
		w.WriteHeader(http.StatusOK)
	} else {
		log.Println("Unknown user tries to log out")
		w.WriteHeader(http.StatusForbidden)
	}
	w.Write([]byte(views.BaseHTML(lc.Get("auth.logout_title"), views.LogoutHTML(can, lc), lc, u)))
}

// handlerLogoutConfirm logs the user out.
//
// TODO: merge into handlerLogout as POST method.
func handlerLogoutConfirm(w http.ResponseWriter, rq *http.Request) {
	user.LogoutFromRequest(w, rq)
	http.Redirect(w, rq, "/", http.StatusSeeOther)
}

// handlerLogin shows the login form.
func handlerLogin(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	if cfg.UseAuth {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
	lc := l18n.FromRequest(rq)
	w.Write([]byte(views.BaseHTML(lc.Get("auth.login_title"), views.LoginHTML(lc), lc, user.EmptyUser())))
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
	if user.HasUsername(username) && user.UserByName(username).Source == "telegram" {
		// Problems is something we put blankets on.
		err = nil
	}

	if !seemsValid {
		err = errors.New("Wrong parameters")
	}

	if err != nil {
		log.Printf("Failed to register ‘%s’ using Telegram: %s", username, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(
			w,
			views.BaseHTML(
				lc.Get("ui.error"),
				fmt.Sprintf(
					`<main class="main-width"><p>%s</p><p>%s</p><p><a href="/login">%s<a></p></main>`,
					lc.Get("auth.error_telegram"),
					err.Error(),
					lc.Get("auth.go_login"),
				),
				lc,
				user.FromRequest(rq),
			),
		)
		return
	}

	errmsg := user.LoginDataHTTP(w, rq, username, "")
	if errmsg != "" {
		log.Printf("Failed to login ‘%s’ using Telegram: %s", username, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(
			w,
			views.BaseHTML(
				"Error",
				fmt.Sprintf(
					`<main class="main-width"><p>%s</p><p>%s</p><p><a href="/login">%s<a></p></main>`,
					lc.Get("auth.error_telegram"),
					err.Error(),
					lc.Get("auth.go_login"),
				),
				lc,
				user.FromRequest(rq),
			),
		)
		return
	}
	log.Printf("Authorize ‘%s’ from Telegram", username)
	http.Redirect(w, rq, "/", http.StatusSeeOther)
}

// handlerLoginData logs the user in.
//
// TODO: merge into handlerLogin as POST method.
func handlerLoginData(w http.ResponseWriter, rq *http.Request) {
	lc := l18n.FromRequest(rq)
	util.PrepareRq(rq)
	var (
		username = util.CanonicalName(rq.PostFormValue("username"))
		password = rq.PostFormValue("password")
		err      = user.LoginDataHTTP(w, rq, username, password)
	)
	if err != "" {
		w.Write([]byte(views.BaseHTML(err, views.LoginErrorHTML(err, lc), lc, user.EmptyUser())))
	} else {
		http.Redirect(w, rq, "/", http.StatusSeeOther)
	}
}
