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
	if util.UseFixedAuth {
		rememberUsers(usersFromFixedCredentials())
	}
	if util.UseRegistration {
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
	users = usersFromFile(util.FixedCredentialsPath, SourceFixed)
	log.Println("Found", len(users), "fixed users")
	return users
}

func usersFromRegistrationCredentials() (users []*User) {
	users = usersFromFile(registrationCredentialsPath(), SourceRegistration)
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
	log.Println("Path for saving tokens:", dir)
	return dir
}

func registrationCredentialsPath() string {
	path := util.RegistrationCredentialsPath
	if path == "" {
		dir, err := xdg.DataFile("mycorrhiza/registration.json")
		if err != nil {
			// No error handling, because the program will fail later anyway when trying to read file ""
			log.Println("Error: cannot get a file to registration credentials, so no registered users will be saved.")
		} else {
			path = dir
		}
	}
	return path
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
	err = ioutil.WriteFile(registrationCredentialsPath(), blob, 0644)
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
		ioutil.WriteFile(tokenStoragePath(), blob, 0644)
	}
}
