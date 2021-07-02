package web

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"
)

func initAuth() {
	if !cfg.UseAuth {
		return
	}
	if cfg.AllowRegistration {
		http.HandleFunc("/register", handlerRegister)
	}
	http.HandleFunc("/login", handlerLogin)
	http.HandleFunc("/login-data", handlerLoginData)
	http.HandleFunc("/logout", handlerLogout)
	http.HandleFunc("/logout-confirm", handlerLogoutConfirm)
}

// handlerRegister both displays the register form (GET) and registers users (POST).
func handlerRegister(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	if !cfg.AllowRegistration {
		w.WriteHeader(http.StatusForbidden)
	}
	if rq.Method == http.MethodGet {
		io.WriteString(
			w,
			views.BaseHTML(
				"Register",
				views.RegisterHTML(rq),
				user.FromRequest(rq),
			),
		)
	} else if rq.Method == http.MethodPost {
		var (
			username = rq.PostFormValue("username")
			password = rq.PostFormValue("password")
			err      = user.Register(username, password)
		)
		if err != nil {
			w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(
				w,
				views.BaseHTML(
					"Register",
					fmt.Sprintf(
						`<main class="main-width"><p>%s</p><p><a href="/register">Try again<a></p></main>`,
						err.Error(),
					),
					user.FromRequest(rq),
				),
			)
		} else {
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
	)
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	if can {
		log.Println("User", u.Name, "tries to log out")
		w.WriteHeader(http.StatusOK)
	} else {
		log.Println("Unknown user tries to log out")
		w.WriteHeader(http.StatusForbidden)
	}
	w.Write([]byte(views.BaseHTML("Logout?", views.LogoutHTML(can), u)))
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
	w.Write([]byte(views.BaseHTML("Login", views.LoginHTML(), user.EmptyUser())))
}

// handlerLoginData logs the user in.
//
// TODO: merge into handlerLogin as POST method.
func handlerLoginData(w http.ResponseWriter, rq *http.Request) {
	util.PrepareRq(rq)
	var (
		username = util.CanonicalName(rq.PostFormValue("username"))
		password = rq.PostFormValue("password")
		err      = user.LoginDataHTTP(w, rq, username, password)
	)
	if err != "" {
		w.Write([]byte(views.BaseHTML(err, views.LoginErrorHTML(err), user.EmptyUser())))
	} else {
		http.Redirect(w, rq, "/", http.StatusSeeOther)
	}
}
