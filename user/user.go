package user

import (
	"net/http"
	"sync"
	"time"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"golang.org/x/crypto/bcrypt"
)

// User is a user.
type User struct {
	// Name is a username. It must follow hypha naming rules.
	Name         string    `json:"name"`
	Group        string    `json:"group"`
	Password     string    `json:"hashed_password"`
	RegisteredAt time.Time `json:"registered_on"`
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

var groups = []string{
	"anon",
	"editor",
	"trusted",
	"telegram",
	"moderator",
	"admin",
}

// Group — Right
var groupRight = map[string]int{
	"anon":      0,
	"editor":    1,
	"trusted":   2,
	"telegram":  2,
	"moderator": 3,
	"admin":     4,
}

func ValidGroup(group string) bool {
	for _, grp := range groups {
		if grp == group {
			return true
		}
	}
	return false
}

func EmptyUser() *User {
	return &User{
		Name:     "anon",
		Group:    "anon",
		Password: "",
	}
}

func (user *User) CanProceed(route string) bool {
	if !cfg.UseAuth {
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

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// ShowLockMaybe redirects to the lock page if the user is anon and the wiki has been configured to use the lock. It returns true if the user was redirected.
func (user *User) ShowLockMaybe(w http.ResponseWriter, rq *http.Request) bool {
	if cfg.Locked && user.Group == "anon" {
		http.Redirect(w, rq, "/lock", http.StatusSeeOther)
		return true
	}
	return false
}
