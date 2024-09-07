package user

import (
	"sort"
	"sync"
)

var users sync.Map
var tokens sync.Map

// YieldUsers creates a channel which iterates existing users.
func YieldUsers() chan *User {
	ch := make(chan *User)
	go func(ch chan *User) {
		users.Range(func(_, v any) bool {
			ch <- v.(*User)
			return true
		})
		close(ch)
	}(ch)
	return ch
}

// ListUsersWithGroup returns a slice with users of desired group.
func ListUsersWithGroup(group string) []string {
	var filtered []string
	for u := range YieldUsers() {
		if u.Group == group {
			filtered = append(filtered, u.Name)
		}
	}
	return filtered
}

// Count returns total users count
func Count() (i uint64) {
	users.Range(func(k, v interface{}) bool {
		i++
		return true
	})
	return i
}

func HasAnyAdmins() bool {
	for u := range YieldUsers() {
		if u.Group == "admin" {
			return true
		}
	}
	return false
}

// HasUsername checks whether the desired user exists
func HasUsername(username string) bool {
	_, has := users.Load(username)
	return has
}

// CredentialsOK checks whether a correct user-password pair is provided
func CredentialsOK(username, password string) bool {
	return ByName(username).isCorrectPassword(password)
}

// ByToken finds a user by provided session token
func ByToken(token string) *User {
	// TODO: Needs more session data -- chekoopa
	if usernameUntyped, ok := tokens.Load(token); ok {
		username := usernameUntyped.(string)
		return ByName(username)
	}
	return EmptyUser()
}

// ByName finds a user by one's username
func ByName(username string) *User {
	if userUntyped, ok := users.Load(username); ok {
		user := userUntyped.(*User)
		return user
	}
	return EmptyUser()
}

// DeleteUser removes a user by one's name and saves user database.
func DeleteUser(name string) error {
	user, loaded := users.LoadAndDelete(name)
	if loaded {
		u := user.(*User)
		u.Name = "anon"
		u.Group = "anon"
		u.Password = ""
		return SaveUserDatabase()
	}
	return nil
}

func commenceSession(username, token string) {
	tokens.Store(token, username)
	dumpTokens()
}

func terminateSession(token string) {
	tokens.Delete(token)
	dumpTokens()
}

func UsersInGroups() (admins []string, moderators []string, editors []string, readers []string) {
	for u := range YieldUsers() {
		switch u.Group {
		// What if we place the users into sorted slices?
		case "admin":
			admins = append(admins, u.Name)
		case "moderator":
			moderators = append(moderators, u.Name)
		case "editor", "trusted":
			editors = append(editors, u.Name)
		case "reader":
			readers = append(readers, u.Name)
		}
	}
	sort.Strings(admins)
	sort.Strings(moderators)
	sort.Strings(editors)
	sort.Strings(readers)
	return
}
