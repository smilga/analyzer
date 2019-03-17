package api

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

var (
	ErrFilterNotFound = errors.New("Error filter not found")
)

type FilterID uuid.UUID

func (id FilterID) MarshalText() (text []byte, err error) {
	return uuid.UUID(id).MarshalText()
}

type FilterStorage interface {
	Save(*Filter) error
	All() ([]*Filter, error)
	Get(FilterID) (*Filter, error)
}

type Filter struct {
	ID          FilterID
	Name        string
	Description string
	Tags        []*Tag
}
