package api

import (
	"errors"
	"net/url"
	"time"
)

// Error definitions
var (
	ErrWebsiteNotFound = errors.New("Error website not found")
)

type WebsiteID int64

type WebsiteStorage interface {
	ByUser(UserID) ([]*Website, error)
	ByFilterID([]FilterID, UserID) ([]*Website, error)
	Get(WebsiteID) (*Website, error)
	Save(*Website) error
	Delete(WebsiteID) error
}

type Website struct {
	ID          WebsiteID  `db:"id"`
	UserID      UserID     `db:"user_id"`
	URL         string     `db:"url"`
	Tags        []*Tag     `db:"-"`
	Matches     []*Match   `db:"-"`
	InspectedAt *time.Time `db:"inspected_at"`
	CreatedAt   *time.Time `db:"created_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

func NewWebsite(uri string, uid UserID) *Website {
	now := time.Now()

	return &Website{
		UserID:      uid,
		URL:         buildURL(uri),
		InspectedAt: nil,
		CreatedAt:   &now,
		DeletedAt:   nil,
	}
}

func buildURL(uri string) string {
	_, err := url.ParseRequestURI(uri)
	if err != nil {
		return uri
	}
	return "http://" + uri
}
