package user

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/util"
	"golang.org/x/crypto/bcrypt"
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
	return UserByToken(cookie.Value)
}

// LogoutFromRequest logs the user in `rq` out and rewrites the cookie in `w`.
func LogoutFromRequest(w http.ResponseWriter, rq *http.Request) {
	cookieFromUser, err := rq.Cookie("mycorrhiza_token")
	if err == nil {
		http.SetCookie(w, cookie("token", "", time.Unix(0, 0)))
		terminateSession(cookieFromUser.Value)
	}
}

// Register registers the given user. If it fails, a non-nil error is returned.
func Register(username, password string) error {
	username = util.CanonicalName(username)
	log.Println("Attempt to register user", username)
	switch {
	case cfg.RegistrationLimit > 0 && Count() >= cfg.RegistrationLimit:
		log.Printf("Limit reached: %d", cfg.RegistrationLimit)
		return fmt.Errorf("Reached the limit of registered users: %d", cfg.RegistrationLimit)
	case HasUsername(username):
		log.Println("Username taken")
		return fmt.Errorf("Username \"%s\" is taken already.", username)
	case !util.IsPossibleUsername(username):
		log.Println("Illegal username:", username)
		return fmt.Errorf("Illegal username \"%s\".", username)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u := User{
		Name:         username,
		Group:        "editor",
		Password:     string(hash),
		RegisteredAt: time.Now(),
	}
	users.Store(username, &u)
	err = SaveUserDatabase()
	if err != nil {
		return err
	}
	return nil
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
