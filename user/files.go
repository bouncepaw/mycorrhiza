package user

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/adrg/xdg"
	"github.com/bouncepaw/mycorrhiza/util"
)

// ReadUsersFromFilesystem reads all user information from filesystem and stores it internally. Call it during initialization.
func ReadUsersFromFilesystem() {
	rememberUsers(usersFromFixedCredentials())
	readTokensToUsers()
}

func usersFromFixedCredentials() []*User {
	var users []*User
	contents, err := ioutil.ReadFile(util.FixedCredentialsPath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(contents, &users)
	if err != nil {
		log.Fatal(err)
	}
	for _, u := range users {
		u.Source = SourceFixed
	}
	log.Println("Found", len(users), "fixed users")
	return users
}

func rememberUsers(uu []*User) {
	// uu is used to not shadow the `users` in `users.go`.
	for _, user := range uu {
		users.Store(user.Name, user)
	}
}

func readTokensToUsers() {
	contents, err := ioutil.ReadFile(tokenStoragePath())
	if os.IsNotExist(err) {
		return
	}
	if err != nil {
		log.Fatal(err)
	}

	var tmp map[string]string
	err = json.Unmarshal(contents, &tmp)
	if err != nil {
		log.Fatal(err)
	}

	for token, username := range tmp {
		commenceSession(username, token)
	}
	log.Println("Found", len(tmp), "active sessions")
}

// Return path to tokens.json. Creates folders if needed.
func tokenStoragePath() string {
	dir, err := xdg.DataFile("mycorrhiza/tokens.json")
	if err != nil {
		log.Fatal(err)
	}
	if strings.HasPrefix(dir, util.WikiDir) {
		log.Fatal("Error: Wiki storage directory includes private config files")
	}
	return dir
}

func dumpTokens() {
	tmp := make(map[string]string)

	tokens.Range(func(k, v interface{}) bool {
		token := k.(string)
		username := v.(string)
		tmp[token] = username
		return true
	})

	blob, err := json.Marshal(tmp)
	if err != nil {
		log.Println(err)
	} else {
		ioutil.WriteFile(tokenStoragePath(), blob, 0644)
	}
}
