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
	ID        TagID      `db:"id"`
	Value     string     `db:"value"`
	CreatedAt *time.Time `db:"created_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func NewTag(value string) *Tag {
	now := time.Now()

	return &Tag{
		Value:     value,
		CreatedAt: &now,
		DeletedAt: nil,
	}
}
