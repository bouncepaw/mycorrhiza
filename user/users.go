package user

import (
	"sync"
)

var AuthUsed bool
var users sync.Map
var tokens sync.Map

func ListUsersWithGroup(group string) []string {
	usersWithTheGroup := []string{}
	users.Range(func(_, v interface{}) bool {
		userobj := v.(*User)

		if userobj.Group == group {
			usersWithTheGroup = append(usersWithTheGroup, userobj.Name)
		}
		return true
	})
	return usersWithTheGroup
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
	return userByName(username).isCorrectPassword(password)
}

func userByToken(token string) *User {
	if usernameUntyped, ok := tokens.Load(token); ok {
		username := usernameUntyped.(string)
		return userByName(username)
	}
	return emptyUser()
}

func userByName(username string) *User {
	if userUntyped, ok := users.Load(username); ok {
		user := userUntyped.(*User)
		return user
	}
	return emptyUser()
}

func commenceSession(username, token string) {
	tokens.Store(token, username)
	go dumpTokens()
}

func terminateSession(token string) {
	tokens.Delete(token)
	go dumpTokens()
}
