package mysql

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/gorm"
	"github.com/smilga/analyzer/api"
)

type ResultStore struct {
	*gorm.DB
}

func (s *ResultStore) Save(result *api.Result) error {
	spew.Dump(result)
	return s.DB.Create(result).Error
}

func NewResultStore(db *gorm.DB) *ResultStore {
	return &ResultStore{db}
}
