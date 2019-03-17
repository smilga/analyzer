package api

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
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

type PatternID uuid.UUID

func (id PatternID) MarshalText() (text []byte, err error) {
	return uuid.UUID(id).MarshalText()
}

func (id *PatternID) UnmarshalJSON(data []byte) error {
	uid := uuid.UUID{}
	err := uid.UnmarshalText(data[1 : (len(data))-1])
	if err != nil {
		return err
	}
	(*id) = PatternID(uid)
	return nil
}

type Pattern struct {
	ID          PatternID
	Type        PatternType
	Value       string
	Description string
	Tags        []*Tag `db:"-"`
	CreatedAt   *time.Time
	DeletedAt   *time.Time
}

type MatchedPattern struct {
	*Pattern
	Match string
}

func NewPattern(t PatternType, v string, d string) *Pattern {
	now := time.Now()

	return &Pattern{
		ID:          PatternID(uuid.NewV4()),
		Type:        t,
		Value:       v,
		Description: d,
		CreatedAt:   &now,
		DeletedAt:   nil,
	}
}
