package user

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/files"
	"github.com/bouncepaw/mycorrhiza/util"
)

// InitUserDatabase loads users, if necessary. Call it during initialization.
func InitUserDatabase() {
	ReadUsersFromFilesystem()
}

// ReadUsersFromFilesystem reads all user information from filesystem and
// stores it internally.
func ReadUsersFromFilesystem() {
	if cfg.UseAuth {
		rememberUsers(usersFromFile())
		readTokensToUsers()
	}
}

func usersFromFile() []*User {
	var users []*User
	contents, err := os.ReadFile(files.UserCredentialsJSON())
	if errors.Is(err, os.ErrNotExist) {
		return users
	}
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(contents, &users)
	if err != nil {
		log.Fatal(err)
	}
	for _, u := range users {
		u.Name = util.CanonicalName(u.Name)
		if u.Source == "" {
			u.Source = "local"
		}
	}
	log.Println("Found", len(users), "users")
	return users
}

func rememberUsers(userList []*User) {
	for _, user := range userList {
		users.Store(user.Name, user)
	}
}

func readTokensToUsers() {
	contents, err := os.ReadFile(files.TokensJSON())
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
		tokens.Store(token, username)
		// commenceSession(username, token)
	}
	log.Println("Found", len(tmp), "active sessions")
}

// SaveUserDatabase stores current user credentials into JSON file by configured path.
func SaveUserDatabase() error {
	return dumpUserCredentials()
}

func dumpUserCredentials() error {
	var userList []*User

	// TODO: lock the map during saving to prevent corruption
	for u := range YieldUsers() {
		userList = append(userList, u)
	}

	blob, err := json.MarshalIndent(userList, "", "\t")
	if err != nil {
		log.Println(err)
		return err
	}

	err = os.WriteFile(files.UserCredentialsJSON(), blob, 0666)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func dumpTokens() {
	tmp := make(map[string]string)

	tokens.Range(func(k, v interface{}) bool {
		token := k.(string)
		username := v.(string)
		tmp[token] = username
		return true
	})

	blob, err := json.MarshalIndent(tmp, "", "\t")
	if err != nil {
		log.Println(err)
		return
	}
	os.WriteFile(files.TokensJSON(), blob, 0666)
}
