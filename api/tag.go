package api

import (
	"errors"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

var (
	ErrTagNotFound = errors.New("Errot tag not found")
)

type TagID uuid.UUID

func (id TagID) MarshalText() (text []byte, err error) {
	return uuid.UUID(id).MarshalText()
}

func (id *TagID) UnmarshalJSON(data []byte) error {
	fmt.Println(data)
	uid := uuid.UUID{}
	err := uid.UnmarshalText(data[1 : (len(data))-1])
	if err != nil {
		return err
	}
	(*id) = TagID(uid)
	return nil
}

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
		ID:        TagID(uuid.NewV4()),
		Value:     value,
		CreatedAt: &now,
		DeletedAt: nil,
	}
}
