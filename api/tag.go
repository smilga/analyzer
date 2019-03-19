package api

import (
	"errors"
	"time"
)

var (
	ErrTagNotFound = errors.New("Errot tag not found")
)

type TagID int64

type TagStorage interface {
	Get(TagID) (*Tag, error)
	All() ([]*Tag, error)
	Save(*Tag) error
}

type Tag struct {
	ID        TagID
	Value     string
	CreatedAt *time.Time
	DeletedAt *time.Time
}

func NewTag(value string) *Tag {
	now := time.Now()

	return &Tag{
		Value:     value,
		CreatedAt: &now,
		DeletedAt: nil,
	}
}
