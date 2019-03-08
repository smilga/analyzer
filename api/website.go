package api

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Error definitions
var (
	ErrWebsiteNotFound = errors.New("Error website not found")
)

type WebsiteStorage interface {
	All(rel bool) ([]*Website, error)
	ByUser(uuid.UUID) ([]*Website, error)
	Get(uuid.UUID) (*Website, error)
	Save(*Website) error
	Delete(uuid.UUID) error
}

type WebsiteID string

type Website struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	URL        string
	Services   []*ServiceIdentity
	CreatedAt  *time.Time
	SearchedAt *time.Time
	DeletedAt  *time.Time
}
