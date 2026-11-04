package userdomain

import (
	"time"
)

// TODO: think about calling it DTO

type UserDTO struct {
	// Name is a username. It must follow hypha naming rules.
	Name         string    `json:"name"`
	Group        GroupKind `json:"group"`
	Password     string    `json:"hashed_password"`
	RegisteredAt time.Time `json:"registered_on"`
	Source       Source    `json:"source"`

	// A note about why HashedPassword is string and not []byte. The reason is
	// simple: golang's json marshals []byte as slice of numbers, which is not
	// acceptable.
}
