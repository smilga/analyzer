package api

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type ResultStorage interface {
	ByWebsite(uuid.UUID) ([]*Result, error)
	LatestByWebsite(uuid.UUID) (*Result, error)
	Save(*Result) error
}

type Match struct {
	PatternID uuid.UUID
	Value     string
}

type DetectedService struct {
	ServiceID uuid.UUID
	Match     []*Match
}

type Result struct {
	Duration         float64
	WebsiteID        uuid.UUID
	CreatedAt        *time.Time
	DetectedServices []*DetectedService
}

func (r *Result) ListServiceIDs() []uuid.UUID {
	ids := make([]uuid.UUID, len(r.DetectedServices))
	for i, s := range r.DetectedServices {
		ids[i] = s.ServiceID
	}
	return ids
}
