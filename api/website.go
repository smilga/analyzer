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
	ByUser(UserID, *Pagination) ([]*Website, int, error)
	ByFilterID([]FilterID, UserID, *Pagination) ([]*Website, int, error)
	Get(WebsiteID) (*Website, error)
	Where(UserID, string, interface{}) ([]*Website, int, error)
	Save(*Website) error
	SaveBatch([]*Website) error
	Delete(WebsiteID) error
	AddTags([]*Website) error
}

type Website struct {
	ID          WebsiteID  `json:"id" db:"id"`
	UserID      UserID     `json:"userId" db:"user_id"`
	URL         string     `json:"url" db:"url"`
	Tags        []*Tag     `json:"tags" db:"-"`
	Matches     []*Match   `json:"matches" db:"-"`
	InspectedAt *time.Time `json:"inspectedAt" db:"inspected_at"`
	CreatedAt   *time.Time `json:"createdAt" db:"created_at"`
	DeletedAt   *time.Time `json:"deletedAt" db:"deleted_at"`
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
