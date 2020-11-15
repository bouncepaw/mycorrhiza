package user

import (
	"log"
	"net/http"
)

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

func CanProceed(rq *http.Request, route string) bool {
	ug := UserAnon
	if u := FromRequest(rq); u != nil {
		ug = u.Group
	}
	return ug.CanAccessRoute(route)
}
