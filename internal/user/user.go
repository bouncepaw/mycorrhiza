package usersvc

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
	"sync"

	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	"github.com/bouncepaw/mycorrhiza/internal/user/domain"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	userdomain.UserDTO
	sync.RWMutex
}

// EmptyUser constructs an anonymous user.
func EmptyUser() *User {
	return &User{
		Name:     "anon",
		Group:    userdomain.GroupAnon,
		Password: "",
		Source:   userdomain.SourceLocal,
	}
}

// WikimindUser constructs the wikimind user, which is to be used for automated wiki edits and has admin privileges.
func WikimindUser() *User {
	return &User{
		Name:     "wikimind",
		Group:    userdomain.GroupAdmin,
		Password: "",
		Source:   userdomain.SourceLocal,
	}
}

// CanProceed checks whether user has rights to visit the provided path (and perform an action).
func (user *User) CanProceed(action userdomain.Action) bool {
	if !cfg.UseAuth {
		return true
	}

	user.RLock()
	defer user.RUnlock()

	return user.Group.Level() >= action.MinimalLevel()
}

func (user *User) isCorrectPassword(password string) bool {
	user.RLock()
	defer user.RUnlock()

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// ShowLockMaybe redirects to the lock page if the user is anon and the wiki has been configured to use the lock. It returns true if the user was redirected.
func (user *User) ShowLockMaybe(w http.ResponseWriter, rq *http.Request) bool {
	if cfg.Locked && user.Group == userdomain.GroupAnon {
		http.Redirect(w, rq, "/lock", http.StatusSeeOther)
		return true
	}
	return false
}

// Sets a new password for the user.
func (user *User) ChangePassword(password string) error {
	if user.Source != userdomain.SourceLocal {
		return fmt.Errorf("Only local users can change their passwords.")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return SaveUserDatabase()
}

// IsValidUsername checks if the given username is valid.
func IsValidUsername(username string) bool {
	for _, r := range username {
		if strings.ContainsRune("?!:#@><*|\"'&%{}/", r) {
			return false
		}
	}
	return username != "anon" &&
		username != "wikimind" &&
		usernameIsWhiteListed(username)
}

func usernameIsWhiteListed(username string) bool {
	if !cfg.UseWhiteList {
		return true
	}
	return slices.Contains(cfg.WhiteList, username)
}
