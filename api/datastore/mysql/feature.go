package mysql

import (
	"github.com/jinzhu/gorm"
	"github.com/smilga/analyzer/api"
)

type FeatureStore struct {
	DB *gorm.DB
}

func (s *FeatureStore) All() ([]*api.Feature, error) {
	fs := []*api.Feature{}
	err := s.DB.Find(&fs).Error
	if err != nil {
		return nil, err
	}
	return fs, nil
}

func (s *FeatureStore) Get(id api.FeatureID) (*api.Feature, error) {
	f := &api.Feature{}

	err := s.DB.First(&f, id).Error
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (s *FeatureStore) Save(f *api.Feature) error {
	return s.DB.Create(f).Error
}

func NewFeatureStore(DB *gorm.DB) *FeatureStore {
	return &FeatureStore{DB}
}
