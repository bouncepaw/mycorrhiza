package user

import (
	"sync"
)

var AuthUsed bool
var users sync.Map
var tokens sync.Map

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
