package main

import (
	"log"
	"net/http"

	"github.com/bouncepaw/mycorrhiza/templates"
	"github.com/bouncepaw/mycorrhiza/user"
)

func init() {
	http.HandleFunc("/login", handlerLogin)
	http.HandleFunc("/login-data", handlerLoginData)
	http.HandleFunc("/logout", handlerLogout)
	http.HandleFunc("/logout-confirm", handlerLogoutConfirm)
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
	w.Write([]byte(base("Logout?", templates.LogoutHTML(can))))
}

func handlerLogoutConfirm(w http.ResponseWriter, rq *http.Request) {
	user.LogoutFromRequest(w, rq)
	http.Redirect(w, rq, "/", http.StatusSeeOther)
}

func handlerLoginData(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	var (
		username = CanonicalName(rq.PostFormValue("username"))
		password = rq.PostFormValue("password")
		err      = user.LoginDataHTTP(w, rq, username, password)
	)
	if err != "" {
		w.Write([]byte(base(err, templates.LoginErrorHTML(err))))
	} else {
		http.Redirect(w, rq, "/", http.StatusSeeOther)
	}
}

func handlerLogin(w http.ResponseWriter, rq *http.Request) {
	log.Println(rq.URL)
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	if user.AuthUsed {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
	w.Write([]byte(base("Login", templates.LoginHTML())))
}
