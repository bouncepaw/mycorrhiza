package web

import (
	"errors"
	"fmt"
	"github.com/bouncepaw/mycorrhiza/viewutil"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"io"
	"log"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/l18n"
	"github.com/bouncepaw/mycorrhiza/user"
	"github.com/bouncepaw/mycorrhiza/util"
	"github.com/bouncepaw/mycorrhiza/views"
	"github.com/gorilla/mux"
)

func initAuth(r *mux.Router) {
	r.HandleFunc("/user-list", handlerUserList)
	r.HandleFunc("/lock", handlerLock)
	// The check below saves a lot of extra checks and lines of codes in other places in this file.
	if !cfg.UseAuth {
		return
	}
	if cfg.AllowRegistration {
		r.HandleFunc("/register", handlerRegister).Methods(http.MethodPost, http.MethodGet)
	}
	if cfg.TelegramEnabled {
		r.HandleFunc("/telegram-login", handlerTelegramLogin)
	}

	if cfg.OidcEnabled {
		r.HandleFunc("/oauth/redirect", handlerOidcRedirect)
		r.HandleFunc("/oauth/login", handlerOidcLogin)
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
			err      = user.LoginDataHTTP(w, username, password)
		)
		if err != nil {
			w.Header().Set("Content-Type", "text/html;charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = io.WriteString(w, views.Base(viewutil.MetaFrom(w, rq), err.Error(), views.LoginError(err.Error(), lc)))
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

	errmsg := user.LoginDataHTTP(w, username, "")
	if errmsg != nil {
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

func handlerOidcLogin(w http.ResponseWriter, rq *http.Request) {
	//lc := l18n.FromRequest(rq)

	_, config, err := util.GetOidcProviderAndConfig(
		cfg.OidcClientId,
		cfg.OidcClientSecret,
		cfg.OidcProviderUri,
		cfg.URL+"/oauth/redirect",
		strings.Split(cfg.OidcPlusSeparatedScopes, "+"),
	)

	if err != nil {
		log.Printf("Failed to log in via OIDC: provider initialization failed: %s", err.Error())
		http.Error(w, "Failed to log in via OIDC: provider initialization failed", http.StatusInternalServerError)
		return
	}

	state, err := util.RandString(16)
	if err != nil {
		// TODO: what is best practice about that?
		log.Printf("Failed to log in via OIDC: RandString behaviour is weird: %s", err.Error())
		http.Error(w, "Failed to log in via OIDC: RandString behaviour is weird", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:  "state",
		Value: state,
		// TODO: there is some weirdness with timezones and cookies in firefox
		MaxAge: 48 * int(time.Hour.Seconds()),
		// TODO: what are these properties below?
		Secure:   rq.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, rq, config.AuthCodeURL(state), http.StatusFound)
}

func handlerOidcRedirect(w http.ResponseWriter, rq *http.Request) {
	lc := l18n.FromRequest(rq)

	provider, config, err := util.GetOidcProviderAndConfig(
		cfg.OidcClientId,
		cfg.OidcClientSecret,
		cfg.OidcProviderUri,
		cfg.URL+"/oauth/redirect",
		strings.Split(cfg.OidcPlusSeparatedScopes, "+"),
	)

	if err != nil {
		log.Printf("Failed to log in via OIDC: provider initialization failed: %s", err.Error())
		http.Error(w, "Failed to log in via OIDC: provider initialization failed", http.StatusInternalServerError)
		return
	}

	// TODO: replace context?
	ctx := context.TODO()

	state, err := rq.Cookie("state")
	if err != nil {
		log.Printf("Failed to log in via OIDC: state cookie not found: %s", err.Error())
		http.Error(w, "Failed to log in via OIDC: state cookie not found", http.StatusInternalServerError)
		return
	}
	if rq.URL.Query().Get("state") != state.Value {
		log.Printf("Failed to log in via OIDC: state cookie is malformed: %s", err.Error())
		http.Error(w, "Failed to log in via OIDC: state cookie is malformed", http.StatusInternalServerError)
		return
	}

	oauth2Token, err := config.Exchange(ctx, rq.URL.Query().Get("code"))
	if err != nil {
		log.Printf("Failed to log in via OIDC: failed to exchange token: %s", err.Error())
		http.Error(w, "Failed to log in via OIDC: failed to exchange token", http.StatusInternalServerError)
		return
	}

	userInfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		log.Printf("Failed to log in via OIDC: failed to get userinfo: %s", err.Error())
		http.Error(w, "Failed to log in via OIDC: failed to get userinfo", http.StatusInternalServerError)
		return
	}

	// TODO: allow usage of different fields. Sometimes there is valid nickname field, so this should be configurable
	username := userInfo.Email
	seemsValid := len(username) > 0 && userInfo.EmailVerified == true

	if !seemsValid {
		err = errors.New("Please set email in your " + cfg.OidcProvider + " instance in order to create account here")
	} else {
		// TODO: Delete this when '@' is allowed in mycomarkup
		username = strings.ReplaceAll(username, "@", "-at-")
		username = util.CanonicalName(username)

		err = user.Register(
			username,
			"", // Password matters not
			"editor",
			// TODO: replace "oidc" with OidcProvider as soon as "source" limitations are removed
			"oidc",
			false,
		)

		// TODO: replace "oidc" with OidcProvider as soon as "source" limitations are removed
		if user.HasUsername(username) && user.ByName(username).Source == "oidc" {
			// Problems is something we put blankets on.
			err = nil
		}
	}

	if err != nil {
		log.Printf("Failed to register ‘%s’ using OIDC: %s", username, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(
			w,
			views.Base(
				viewutil.MetaFrom(w, rq),
				lc.Get("ui.error"),
				fmt.Sprintf(
					`<main class="main-width"><p>%s</p><p>%s</p><p><a href="/login">%s<a></p></main>`,
					lc.Get("auth.error_oidc"),
					err.Error(),
					lc.Get("auth.go_login"),
				),
			),
		)
		return
	}

	errmsg := user.LoginDataHTTP(w, username, "")
	if errmsg != nil {
		log.Printf("Failed to login ‘%s’ using OIDC: %s", username, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(
			w,
			views.Base(
				viewutil.MetaFrom(w, rq),
				"Error",
				fmt.Sprintf(
					`<main class="main-width"><p>%s</p><p>%s</p><p><a href="/login">%s<a></p></main>`,
					lc.Get("auth.error_oidc"),
					err.Error(),
					lc.Get("auth.go_login"),
				),
			),
		)
		return
	}
	log.Printf("Authorize ‘%s’ from OIDC", username)
	http.Redirect(w, rq, "/", http.StatusSeeOther)
}
