package api

import (
	"errors"
	"net/url"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Error definitions
var (
	ErrWebsiteNotFound = errors.New("Error website not found")
)

type WebsiteID uuid.UUID

func (id WebsiteID) MarshalText() (text []byte, err error) {
	return uuid.UUID(id).MarshalText()
}

type WebsiteStorage interface {
	All(rel bool) ([]*Website, error)
	ByUser(UserID) ([]*Website, error)
	Get(WebsiteID) (*Website, error)
	Save(*Website) error
	Delete(WebsiteID) error
}

type Website struct {
	ID              WebsiteID
	UserID          UserID
	URL             string
	MatchedPatterns []*MatchedPattern `db:"-"`
	SearchedAt      *time.Time
	CreatedAt       *time.Time
	DeletedAt       *time.Time
}

func NewWebsite(uri string, uid UserID) *Website {
	now := time.Now()

	return &Website{
		ID:         WebsiteID(uuid.NewV4()),
		UserID:     uid,
		URL:        buildURL(uri),
		SearchedAt: nil,
		CreatedAt:  &now,
		DeletedAt:  nil,
	}
}

func buildURL(uri string) string {
	_, err := url.ParseRequestURI(uri)
	if err != nil {
		return uri
	}
	return "http://" + uri
}
