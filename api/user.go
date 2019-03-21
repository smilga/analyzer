package api

import (
	"errors"
	"time"
)

// Error definitions
var (
	ErrUserNotFound = errors.New("Error user not found")
)

type UserID int64

type UserStorage interface {
	Save(*User) error
	ByEmail(string) (*User, error)
	ByID(UserID) (*User, error)
}

type User struct {
	ID        UserID      `db:"id"`
	Name      string      `db:"name"`
	Email     string      `db:"email"`
	Password  Cryptstring `db:"password" json:"-"`
	CreatedAt *time.Time  `db:"created_at"`
	UpdatedAt *time.Time  `db:"updated_at"`
	DeletedAt *time.Time  `db:"deleted_at"`
}
