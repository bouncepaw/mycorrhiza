package main

import (
	"github.com/bouncepaw/mycorrhiza/cfg"
	"io"
	"log"
	"net/http"

	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"
)

func init() {
	http.HandleFunc("/register", handlerRegister)
	http.HandleFunc("/login", handlerLogin)
	http.HandleFunc("/login-data", handlerLoginData)
	http.HandleFunc("/logout", handlerLogout)
	http.HandleFunc("/logout-confirm", handlerLogoutConfirm)
}

func handlerRegister(w http.ResponseWriter, rq *http.Request) {
	prepareRq(rq)
	if !cfg.UseRegistration {
		w.WriteHeader(http.StatusForbidden)
	}
	if rq.Method == http.MethodGet {
		io.WriteString(
			w,
			base(
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
			io.WriteString(w, err.Error())
		} else {
			user.LoginDataHTTP(w, rq, username, password)
			http.Redirect(w, rq, "/"+rq.URL.RawQuery, http.StatusSeeOther)
		}
	}
}

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
	w.Write([]byte(base("Logout?", views.LogoutHTML(can), u)))
}

func handlerLogoutConfirm(w http.ResponseWriter, rq *http.Request) {
	user.LogoutFromRequest(w, rq)
	http.Redirect(w, rq, "/", http.StatusSeeOther)
}

func handlerLoginData(w http.ResponseWriter, rq *http.Request) {
	prepareRq(rq)
	var (
		username = util.CanonicalName(rq.PostFormValue("username"))
		password = rq.PostFormValue("password")
		err      = user.LoginDataHTTP(w, rq, username, password)
	)
	if err != "" {
		w.Write([]byte(base(err, views.LoginErrorHTML(err), user.EmptyUser())))
	} else {
		http.Redirect(w, rq, "/", http.StatusSeeOther)
	}
}

func handlerLogin(w http.ResponseWriter, rq *http.Request) {
	prepareRq(rq)
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	if user.AuthUsed {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
	w.Write([]byte(base("Login", views.LoginHTML(), user.EmptyUser())))
}
