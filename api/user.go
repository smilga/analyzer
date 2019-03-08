package api

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type UserStorage interface {
	ByEmail(string) (*User, error)
	ByID(uuid.UUID) (*User, error)
}

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  Cryptstring `json:"-"`
	CreatedAt *time.Time
}
