package mysql

import (
	"github.com/jinzhu/gorm"
	"github.com/smilga/analyzer/api"
)

type UserStore struct {
	DB *gorm.DB
}

func (s *UserStore) All() error {
	us := []*api.User{}
	return s.DB.Find(&us).Error
}

func (s *UserStore) Save(u *api.User) error {
	return s.DB.Create(u).Error
}

func (s *UserStore) ByID(id api.UserID) (*api.User, error) {
	var u api.User
	err := s.DB.Where("id = ?", int(id)).First(&u).Error
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *UserStore) ByEmail(email string) (*api.User, error) {
	var u api.User
	err := s.DB.Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func NewUserStore(DB *gorm.DB) *UserStore {
	return &UserStore{DB}
}
