package user

import (
	"sync"

	"golang.org/x/crypto/bcrypt"
)

// UserSource shows where is the user data gotten from.
type UserSource int

const (
	SourceUnknown UserSource = iota
	// SourceFixed is used with users that are predefined using fixed auth
	SourceFixed
	// SourceRegistration is used with users that are registered through the register form
	SourceRegistration
)

// User is a user.
type User struct {
	// Name is a username. It must follow hypha naming rules.
	Name           string     `json:"name"`
	Group          string     `json:"group"`
	Password       string     `json:"password"`        // for fixed
	HashedPassword string     `json:"hashed_password"` // for registered
	Source         UserSource `json:"-"`
	sync.RWMutex

	// A note about why HashedPassword is string and not []byte. The reason is
	// simple: golang's json marshals []byte as slice of numbers, which is not
	// acceptable.
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
	"admin":               4,
	"admin/shutdown":      4,
}

// Group — Right
var groupRight = map[string]int{
	"anon":      0,
	"editor":    1,
	"trusted":   2,
	"moderator": 3,
	"admin":     4,
}

func EmptyUser() *User {
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

	switch user.Source {
	case SourceFixed:
		return password == user.Password
	case SourceRegistration:
		err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
		return err == nil
	}
	return false
}
