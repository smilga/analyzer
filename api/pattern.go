package api

import (
	"errors"
	"time"
)

var (
	ErrPatternNotFound = errors.New("Error pattern not found")
)

// Possible pattern types
const (
	JSSource PatternType = "js_source"
	HTML     PatternType = "html"
	Resource PatternType = "resource"
)

type PatternStorage interface {
	Save(*Pattern) error
	Delete(PatternID) error
	Get(PatternID) (*Pattern, error)
	All() ([]*Pattern, error)
}

type PatternType string

type PatternID int64

type Pattern struct {
	ID          PatternID   `db:"id"`
	Type        PatternType `db:"type"`
	Value       string      `db:"value"`
	Description string      `db:"description"`
	Tags        []*Tag      `db:"-"`
	CreatedAt   *time.Time  `db:"created_at"`
	UpdatedAt   *time.Time  `db:"updated_at"`
	DeletedAt   *time.Time  `db:"deleted_at"`
}

func NewPattern(t PatternType, v string, d string) *Pattern {
	now := time.Now()

	return &Pattern{
		Type:        t,
		Value:       v,
		Description: d,
		CreatedAt:   &now,
		DeletedAt:   nil,
	}
}
