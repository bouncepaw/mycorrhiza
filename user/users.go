package user

import (
	"sync"
)

var AuthUsed bool
var users sync.Map
var tokens sync.Map

func YieldUsers() chan *User {
	ch := make(chan *User)
	go func(ch chan *User) {
		users.Range(func(_, v interface{}) bool {
			ch <- v.(*User)
			return true
		})
		close(ch)
	}(ch)
	return ch
}

func ListUsersWithGroup(group string) []string {
	usersWithTheGroup := []string{}
	for u := range YieldUsers() {
		if u.Group == group {
			usersWithTheGroup = append(usersWithTheGroup, u.Name)
		}
	}
	return usersWithTheGroup
}

func CountRegistered() int {
	i := 0
	users.Range(func(k, v interface{}) bool {
		if v.(*User).Source == SourceRegistration {
			i++
		}
		return true
	})
	return i
}

func Count() int {
	i := 0
	users.Range(func(k, v interface{}) bool {
		i++
		return true
	})
	return i
}

func HasUsername(username string) bool {
	_, has := users.Load(username)
	return has
}

func CredentialsOK(username, password string) bool {
	return UserByName(username).isCorrectPassword(password)
}

func UserByToken(token string) *User {
	if usernameUntyped, ok := tokens.Load(token); ok {
		username := usernameUntyped.(string)
		return UserByName(username)
	}
	return EmptyUser()
}

func UserByName(username string) *User {
	if userUntyped, ok := users.Load(username); ok {
		user := userUntyped.(*User)
		return user
	}
	return EmptyUser()
}

func commenceSession(username, token string) {
	tokens.Store(token, username)
	dumpTokens()
}

func terminateSession(token string) {
	tokens.Delete(token)
	dumpTokens()
}
