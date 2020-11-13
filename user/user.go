package user

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/bouncepaw/mycorrhiza/util"
)

type FixedUserStorage struct {
	Users []*User
}

var UserStorage = FixedUserStorage{}

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
