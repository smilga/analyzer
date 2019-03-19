package api

import (
	"time"
)

type UserID int64

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
