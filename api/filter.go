package api

import (
	"errors"
)

var (
	ErrFilterNotFound = errors.New("Error filter not found")
)

type FilterID int64

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
