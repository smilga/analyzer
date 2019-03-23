package api

import (
	"errors"
	"time"
)

var (
	ErrFilterNotFound = errors.New("Error filter not found")
)

type FilterID int64

type FilterStorage interface {
	Save(*Filter) error
	Get(FilterID) (*Filter, error)
	All() ([]*Filter, error)
}

type Filter struct {
	ID          FilterID   `db:"id"`
	Name        string     `db:"name"`
	Description string     `db:"description"`
	Tags        []*Tag     `db:"-"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}
