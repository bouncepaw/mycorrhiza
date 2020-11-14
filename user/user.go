package user

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/bouncepaw/mycorrhiza/util"
)

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
		UserStorage.Tokens[token] = username
		log.Println("New token for", username, "is", token)
	}
	return token, err
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
	Tokens map[string]string
}

var UserStorage = FixedUserStorage{Tokens: make(map[string]string)}

func PopulateFixedUserStorage() {
	contents, err := ioutil.ReadFile(util.FixedCredentialsPath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(contents, &UserStorage.Users)
	if err != nil {
		log.Fatal(err)
	}
	for _, user := range UserStorage.Users {
		user.Group = groupFromString(user.GroupString)
	}
	log.Println("Found", len(UserStorage.Users), "fixed users")
}

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

func groupFromString(s string) UserGroup {
	switch s {
	case "admin":
		return UserAdmin
	case "moderator":
		return UserModerator
	case "trusted":
		return UserTrusted
	case "editor":
		return UserEditor
	default:
		log.Fatal("Unknown user group", s)
		return UserAnon
	}
}

// UserGroup represents a group that a user is part of.
type UserGroup int

const (
	// UserAnon is the default user group which all unauthorized visitors have.
	UserAnon UserGroup = iota
	// UserEditor is a user who can edit and upload stuff.
	UserEditor
	// UserTrusted is a trusted editor who can also rename stuff.
	UserTrusted
	// UserModerator is a moderator who can also delete stuff.
	UserModerator
	// UserAdmin can do everything.
	UserAdmin
)

var minimalRights = map[string]UserGroup{
	"edit":           UserEditor,
	"upload-binary":  UserEditor,
	"upload-text":    UserEditor,
	"rename-ask":     UserTrusted,
	"rename-confirm": UserTrusted,
	"delete-ask":     UserModerator,
	"delete-confirm": UserModerator,
	"reindex":        UserAdmin,
}

func (ug UserGroup) CanAccessRoute(route string) bool {
	if !AuthUsed {
		return true
	}
	if minimalRight, ok := minimalRights[route]; ok {
		if ug >= minimalRight {
			return true
		}
		return false
	}
	return true
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
