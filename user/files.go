package user

import (
	"encoding/json"
	"github.com/bouncepaw/mycorrhiza/cfg"
	"io/ioutil"
	"log"
	"os"

	"github.com/bouncepaw/mycorrhiza/files"
)

// InitUserDatabase checks the configuration for auth methods and loads users
// if necessary. Call it during initialization.
func InitUserDatabase() {
	AuthUsed = cfg.UseFixedAuth || cfg.UseRegistration

	if AuthUsed && (cfg.FixedAuthCredentialsPath != "" || cfg.RegistrationCredentialsPath != "") {
		ReadUsersFromFilesystem()
	}
}

// ReadUsersFromFilesystem reads all user information from filesystem and stores it internally.
func ReadUsersFromFilesystem() {
	if cfg.UseFixedAuth {
		rememberUsers(usersFromFixedCredentials())
	}
	if cfg.UseRegistration {
		rememberUsers(usersFromRegistrationCredentials())
	}
	readTokensToUsers()
}

func usersFromFile(path string, source UserSource) (users []*User) {
	contents, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(contents, &users)
	if err != nil {
		log.Fatal(err)
	}
	for _, u := range users {
		u.Source = source
	}
	return users
}

func usersFromFixedCredentials() (users []*User) {
	users = usersFromFile(files.FixedCredentialsJSON(), SourceFixed)
	log.Println("Found", len(users), "fixed users")
	return users
}

func usersFromRegistrationCredentials() (users []*User) {
	users = usersFromFile(files.RegistrationCredentialsJSON(), SourceRegistration)
	log.Println("Found", len(users), "registered users")
	return users
}

func rememberUsers(uu []*User) {
	// uu is used to not shadow the `users` in `users.go`.
	for _, user := range uu {
		users.Store(user.Name, user)
	}
}

func readTokensToUsers() {
	contents, err := ioutil.ReadFile(files.TokensJSON())
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

func dumpRegistrationCredentials() error {
	tmp := []*User{}

	for u := range YieldUsers() {
		if u.Source != SourceRegistration {
			continue
		}
		copiedUser := u
		copiedUser.Password = ""
		tmp = append(tmp, copiedUser)
	}

	blob, err := json.Marshal(tmp)
	if err != nil {
		log.Println(err)
		return err
	}
	err = ioutil.WriteFile(files.RegistrationCredentialsJSON(), blob, 0644)
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

	blob, err := json.Marshal(tmp)
	if err != nil {
		log.Println(err)
	} else {
		ioutil.WriteFile(files.TokensJSON(), blob, 0644)
	}
}
