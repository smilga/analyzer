package inmem

import (
	"errors"
	"log"

	uuid "github.com/satori/go.uuid"
	"github.com/smilga/analyzer/api"
	"golang.org/x/crypto/bcrypt"
)

// Error definitions
var (
	ErrUserNotFound = errors.New("Error user not found")
)

var users = []*api.User{
	&api.User{
		ID:       api.UserID(uuid.Must(uuid.FromString("3fba8a7b-274c-4613-a7a8-1cae01ce8a98"))),
		Name:     "Kaspars Smilga",
		Email:    "smilga.kaspars@gmail.com",
		Password: cryptPass("pass"),
	},
	&api.User{
		ID:       api.UserID(uuid.Must(uuid.FromString("00311786-2151-4b9a-bb3a-45e7227886f6"))),
		Name:     "Admin",
		Email:    "admin@inspected.tech",
		Password: cryptPass("pass"),
	},
}

func cryptPass(s string) api.Cryptstring {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	return api.Cryptstring(hash)
}

type UserStore struct {
	users []*api.User
}

func (s *UserStore) ByEmail(email string) (*api.User, error) {
	for _, u := range s.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, ErrUserNotFound
}

func (s *UserStore) ByID(uid api.UserID) (*api.User, error) {
	for _, u := range s.users {
		if u.ID == uid {
			return u, nil
		}
	}
	return nil, ErrUserNotFound
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: users,
	}
}
