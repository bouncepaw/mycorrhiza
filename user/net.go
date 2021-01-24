package user

import (
	"log"
	"net/http"
	"time"

	"github.com/bouncepaw/mycorrhiza/util"
)

// CanProceed returns `true` if the user in `rq` has enough rights to access `route`.
func CanProceed(rq *http.Request, route string) bool {
	return FromRequest(rq).CanProceed(route)
}

// FromRequest returns user from `rq`. If there is no user, an anon user is returned instead.
func FromRequest(rq *http.Request) *User {
	cookie, err := rq.Cookie("mycorrhiza_token")
	if err != nil {
		return EmptyUser()
	}
	return userByToken(cookie.Value)
}

// LogoutFromRequest logs the user in `rq` out and rewrites the cookie in `w`.
func LogoutFromRequest(w http.ResponseWriter, rq *http.Request) {
	cookieFromUser, err := rq.Cookie("mycorrhiza_token")
	if err == nil {
		http.SetCookie(w, cookie("token", "", time.Unix(0, 0)))
		terminateSession(cookieFromUser.Value)
	}
}

// LoginDataHTTP logs such user in and returns string representation of an error if there is any.
func LoginDataHTTP(w http.ResponseWriter, rq *http.Request, username, password string) string {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	if !HasUsername(username) {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Unknown username", username, "was entered")
		return "unknown username"
	}
	if !CredentialsOK(username, password) {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("A wrong password was entered for username", username)
		return "wrong password"
	}
	token, err := AddSession(username)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return err.Error()
	}
	http.SetCookie(w, cookie("token", token, time.Now().Add(365*24*time.Hour)))
	return ""
}

// AddSession saves a session for `username` and returns a token to use.
func AddSession(username string) (string, error) {
	token, err := util.RandomString(16)
	if err == nil {
		commenceSession(username, token)
		log.Println("New token for", username, "is", token)
	}
	return token, err
}

// A handy cookie constructor
func cookie(name_suffix, val string, t time.Time) *http.Cookie {
	return &http.Cookie{
		Name:    "mycorrhiza_" + name_suffix,
		Value:   val,
		Expires: t,
		Path:    "/",
	}
}
