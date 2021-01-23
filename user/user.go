package user

import (
	"sync"
)

// User is a user.
type User struct {
	// Name is a username. It must follow hypha naming rules.
	Name     string `json:"name"`
	Group    string `json:"group"`
	Password string `json:"password"`
	sync.RWMutex
}

// Route — Right (more is more right)
var minimalRights = map[string]int{
	"edit":                1,
	"upload-binary":       1,
	"upload-text":         1,
	"rename-ask":          2,
	"rename-confirm":      2,
	"unattach-ask":        2,
	"unattach-confirm":    2,
	"update-header-links": 3,
	"delete-ask":          3,
	"delete-confirm":      3,
	"reindex":             4,
}

// Group — Right
var groupRight = map[string]int{
	"anon":      0,
	"editor":    1,
	"trusted":   2,
	"moderator": 3,
	"admin":     4,
}

func emptyUser() *User {
	return &User{
		Name:     "anon",
		Group:    "anon",
		Password: "",
	}
}

func (user *User) CanProceed(route string) bool {
	if !AuthUsed {
		return true
	}

	user.RLock()
	defer user.RUnlock()

	right, _ := groupRight[user.Group]
	minimalRight, _ := minimalRights[route]
	if right >= minimalRight {
		return true
	}
	return false
}

func (user *User) isCorrectPassword(password string) bool {
	user.RLock()
	defer user.RUnlock()

	return password == user.Password
}
