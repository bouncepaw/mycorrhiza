package user

import (
	"log"
	"net/http"
	"time"

	"github.com/bouncepaw/mycorrhiza/util"
)

func (u *User) OrAnon() *User {
	if u == nil {
		return &User{}
	}
	return u
}

func LogoutFromRequest(w http.ResponseWriter, rq *http.Request) {
	cookieFromUser, err := rq.Cookie("mycorrhiza_token")
	if err == nil {
		http.SetCookie(w, cookie("token", "", time.Unix(0, 0)))
		terminateSession(cookieFromUser.Value)
	}
}

func (us *FixedUserStorage) userByToken(token string) *User {
	if user, ok := us.Tokens[token]; ok {
		return user
	}
	return nil
}

func (us *FixedUserStorage) userByName(username string) *User {
	for _, user := range us.Users {
		if user.Name == username {
			return user
		}
	}
	return nil
}

func FromRequest(rq *http.Request) *User {
	cookie, err := rq.Cookie("mycorrhiza_token")
	if err != nil {
		return nil
	}
	return UserStorage.userByToken(cookie.Value)
}

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
	http.SetCookie(w, cookie("token", token, time.Now().Add(14*24*time.Hour)))
	return ""
}

// AddSession saves a session for `username` and returns a token to use.
func AddSession(username string) (string, error) {
	token, err := util.RandomString(16)
	if err == nil {
		for _, user := range UserStorage.Users {
			if user.Name == username {
				UserStorage.Tokens[token] = user
				go dumpTokens()
			}
		}
		log.Println("New token for", username, "is", token)
	}
	return token, err
}

func terminateSession(token string) {
	delete(UserStorage.Tokens, token)
	go dumpTokens()
}

func HasUsername(username string) bool {
	for _, user := range UserStorage.Users {
		if user.Name == username {
			return true
		}
	}
	return false
}

func CredentialsOK(username, password string) bool {
	for _, user := range UserStorage.Users {
		if user.Name == username && user.Password == password {
			return true
		}
	}
	return false
}

type FixedUserStorage struct {
	Users  []*User
	Tokens map[string]*User
}

var UserStorage = FixedUserStorage{Tokens: make(map[string]*User)}

// AuthUsed shows if a method of authentication is used. You should set it by yourself.
var AuthUsed bool

// User is a user.
type User struct {
	// Name is a username. It must follow hypha naming rules.
	Name string `json:"name"`
	// Group the user is part of.
	Group       UserGroup `json:"-"`
	GroupString string    `json:"group"`
	Password    string    `json:"password"`
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
