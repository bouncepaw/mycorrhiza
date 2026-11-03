// Package storage persists user accounts and sessions as JSON files.
package storage

import (
	"github.com/bouncepaw/mycorrhiza/internal/user/domain"
	"github.com/bouncepaw/mycorrhiza/util/bxjson"
)

type Storage struct {
	users  *bxjson.JSONFile[[]userdomain.UserDTO]
	tokens *bxjson.JSONFile[map[string]string]
}

func New(usersPath, tokensPath string) *Storage {
	return &Storage{
		users:  bxjson.New[[]userdomain.UserDTO](usersPath),
		tokens: bxjson.New[map[string]string](tokensPath),
	}
}

func (s *Storage) LoadUsers() ([]userdomain.UserDTO, error) {
	return s.users.Read()
}

func (s *Storage) SaveUsers(users []userdomain.UserDTO) error {
	return s.users.Write(users)
}

func (s *Storage) LoadTokens() (map[string]string, error) {
	return s.tokens.Read()
}

func (s *Storage) SaveTokens(tokens map[string]string) error {
	return s.tokens.Write(tokens)
}
