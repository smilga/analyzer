package api

import (
	"errors"
	"fmt"
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

type Pagination struct {
	limit  int
	page   int
	search string
}

func (p *Pagination) Offset() int {
	return (p.Page() - 1) * p.limit
}

func (p *Pagination) Limit() int {
	if p.limit == 0 {
		return 2147483647
	}
	return p.limit
}

func (p *Pagination) Page() int {
	if p.page == 0 {
		return 1
	}
	return p.page
}

func (p *Pagination) Search() string {
	return fmt.Sprintf("%%%s%%", p.search)
}

func NewPagination(limit int, page int) *Pagination {
	if limit == 0 {
		limit = 10
	}
	if page == 0 {
		page = 1
	}
	return &Pagination{
		limit,
		page,
		"",
	}
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
