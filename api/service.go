package api

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Error definitions
var (
	ErrServiceNotFound = errors.New("Error service not found")
)

type ServiceStorage interface {
	All(rel bool) ([]*Service, error)
	ByUser(uuid.UUID) ([]*Service, error)
	ManyByUser(uuid.UUID, []uuid.UUID) ([]*Service, error)
	Get(uuid.UUID) (*Service, error)
	Save(*Service) error
	Delete(uuid.UUID) error
}

type ServiceIdentity struct {
	ID        uuid.UUID
	Name      string
	LogoURL   string
	CratedAt  *time.Time
	DeletedAt *time.Time
}

type Service struct {
	*ServiceIdentity
	UserID   uuid.UUID
	Patterns []*Pattern
}

func (s *Service) Pattern(id uuid.UUID) (*Pattern, error) {
	for _, p := range s.Patterns {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, errors.New("Pattern not found")
}
