package user

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"os"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/files"
	"github.com/bouncepaw/mycorrhiza/util"
)

// InitUserDatabase checks the configuration for auth methods and loads users
// if necessary. Call it during initialization.
func InitUserDatabase() {
	AuthUsed = cfg.UseFixedAuth || cfg.UseRegistration

	if AuthUsed {
		ReadUsersFromFilesystem()
	}
}

// ReadUsersFromFilesystem reads all user information from filesystem and stores it internally.
func ReadUsersFromFilesystem() {
	if cfg.UseFixedAuth {
		// This one will be removed.
		rememberUsers(usersFromFixedCredentials())
	}

	// And this one will be renamed to just "users" in the future.
	rememberUsers(usersFromRegistrationCredentials())

	// Migrate fixed users to registered
	tryToMigrate()

	readTokensToUsers()
}

func tryToMigrate() {
	// Fixed authorization should be removed by the next release (1.13).
	// So let's try to help fixed users and migrate them over!

	migrated := 0

	for user := range YieldUsers() {
		if user.Source == SourceFixed {
			hashedPasswd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				log.Fatal("Failed to migrate fixed users:", err)
			}

			user.Password = ""
			user.HashedPassword = string(hashedPasswd)
			user.Source = SourceRegistration
			migrated++
		}
	}

	if migrated > 0 {
		if err := dumpRegistrationCredentials(); err != nil {
			log.Fatal("Failed to migrate fixed users:", err)
		}
		log.Printf("Migrated %d users", migrated)
	}
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
		u.Name = util.CanonicalName(u.Name)
		u.Source = source
	}
	return users
}

func usersFromFixedCredentials() []*User {
	users := usersFromFile(files.FixedCredentialsJSON(), SourceFixed)
	log.Println("Found", len(users), "fixed users")
	return users
}

func usersFromRegistrationCredentials() []*User {
	users := usersFromFile(files.RegistrationCredentialsJSON(), SourceRegistration)
	log.Println("Found", len(users), "registered users")
	return users
}

func rememberUsers(userList []*User) {
	for _, user := range userList {
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

	blob, err := json.MarshalIndent(tmp, "", "\t")
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

	blob, err := json.MarshalIndent(tmp, "", "\t")
	if err != nil {
		log.Println(err)
	} else {
		ioutil.WriteFile(files.TokensJSON(), blob, 0644)
	}
}
