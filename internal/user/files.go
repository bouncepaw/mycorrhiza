package usersvc

import (
	"log/slog"
	"os"

	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	"github.com/bouncepaw/mycorrhiza/internal/user/domain"
	"github.com/bouncepaw/mycorrhiza/internal/user/storage"
	"github.com/bouncepaw/mycorrhiza/util"
)

var database *storage.Storage

// InitUserDatabase loads users, if necessary. Call it during initialization.
func InitUserDatabase(store *storage.Storage) {
	database = store
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
	records, err := database.LoadUsers()
	if err != nil {
		slog.Error("Failed to load users", "err", err)
		os.Exit(1)
	}
	users := make([]*User, 0, len(records))
	for _, record := range records {
		users = append(users, &User{UserDTO: record})
	}

	migrated := false
	for _, u := range users {
		// Upgrade accounts from the legacy trusted editor group.
		if u.Group == "trusted" {
			u.Group = userdomain.GroupEditor
			migrated = true
		}
		u.Name = util.CanonicalName(u.Name)
		if u.Source == "" {
			u.Source = userdomain.SourceLocal
		}
	}
	if migrated {
		slog.Info("Migrated from the legacy trusted editors group")
		if err := writeUserCredentials(users); err != nil {
			slog.Error("Failed to persist trusted editor migration", "err", err)
		}
	}
	slog.Info("Indexed users", "n", len(users))
	return users
}

func rememberUsers(userList []*User) {
	for _, user := range userList {
		if !IsValidUsername(user.Name) {
			continue
		}
		users.Store(user.Name, user)
	}
}

func readTokensToUsers() {
	tmp, err := database.LoadTokens()
	if err != nil {
		slog.Error("Failed to load sessions", "err", err)
		os.Exit(1)
	}

	for token, username := range tmp {
		tokens.Store(token, username)
	}
	slog.Info("Indexed active sessions", "n", len(tmp))
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

	return writeUserCredentials(userList)
}

func writeUserCredentials(userList []*User) error {
	records := make([]userdomain.UserDTO, 0, len(userList))
	for _, u := range userList {
		u.RLock()
		records = append(records, u.UserDTO)
		u.RUnlock()
	}
	return database.SaveUsers(records)
}

func dumpTokens() {
	tmp := make(map[string]string)

	tokens.Range(func(k, v any) bool {
		token := k.(string)
		username := v.(string)
		tmp[token] = username
		return true
	})

	if err := database.SaveTokens(tmp); err != nil {
		slog.Error("Failed to save sessions", "err", err)
	}
}
