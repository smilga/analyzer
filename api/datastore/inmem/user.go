package inmem

import (
	"log"

	"github.com/smilga/analyzer/api"
	"golang.org/x/crypto/bcrypt"
)

var users = []*api.User{
	&api.User{
		ID:       1,
		Name:     "Kaspars Smilga",
		Email:    "smilga.kaspars@gmail.com",
		Password: cryptPass("pass"),
	},
	&api.User{
		ID:       2,
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

func (s *UserStore) Save(target *api.User) error {
	if target.ID == 0 {
		var last int64
		for _, n := range s.users {
			if int64(n.ID) > last {
				last = int64(n.ID)
			}
		}
		target.ID = api.UserID(last + 1)
	}

	for i, user := range s.users {
		if user.ID == target.ID {
			s.users = append(s.users[:i], s.users[i+1:]...)
		}
	}
	s.users = append(s.users, target)

	return nil
}

func (s *UserStore) ByEmail(email string) (*api.User, error) {
	for _, u := range s.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, api.ErrUserNotFound
}

func (s *UserStore) ByID(uid api.UserID) (*api.User, error) {
	for _, u := range s.users {
		if u.ID == uid {
			return u, nil
		}
	}
	return nil, api.ErrUserNotFound
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: users,
	}
}
