package mysql

import (
	"github.com/jinzhu/gorm"
	"github.com/smilga/analyzer/api"
)

type TagStore struct {
	DB *gorm.DB
}

func (s *TagStore) All() ([]*api.Tag, error) {
	ts := []*api.Tag{}
	err := s.DB.Find(&ts).Error
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func (s *TagStore) Get(id api.TagID) (*api.Tag, error) {
	t := &api.Tag{}

	err := s.DB.First(&t, id).Error
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TagStore) Save(t *api.Tag) error {
	return s.DB.Create(t).Error
}

func NewTagStore(DB *gorm.DB) *TagStore {
	return &TagStore{DB}
}
