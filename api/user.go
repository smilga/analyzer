package api

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type UserID uuid.UUID

func (id UserID) String() string {
	return uuid.UUID(id).String()
}
func (id UserID) MarshalText() (text []byte, err error) {
	return uuid.UUID(id).MarshalText()
}

type UserStorage interface {
	ByEmail(string) (*User, error)
	ByID(UserID) (*User, error)
}

type User struct {
	ID        UserID
	Name      string
	Email     string
	Password  Cryptstring `json:"-"`
	CreatedAt *time.Time
	DeletedAt *time.Time
}
