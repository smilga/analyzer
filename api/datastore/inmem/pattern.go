package inmem

import (
	uuid "github.com/satori/go.uuid"
	"github.com/smilga/analyzer/api"
)

type PatternStore struct {
	patterns []*api.Pattern
}

func (s *PatternStore) Save(target *api.Pattern) error {
	for i, p := range s.patterns {
		if p.ID == api.PatternID(target.ID) {
			s.patterns = append(s.patterns[:i], s.patterns[i+1:]...)
		}
	}
	s.patterns = append(s.patterns, target)

	return nil
}

func (s *PatternStore) All() ([]*api.Pattern, error) {
	return s.patterns, nil
}

func (s *PatternStore) Get(id api.PatternID) (*api.Pattern, error) {
	for _, p := range s.patterns {
		if p.ID == api.PatternID(id) {
			return p, nil
		}
	}
	return nil, api.ErrPatternNotFound
}

func (s *PatternStore) Delete(id api.PatternID) error {
	for i, pattern := range s.patterns {
		if pattern.ID == id {
			s.patterns = append(s.patterns[:i], s.patterns[i+1:]...)
		}
	}
	return nil
}

func NewPatternStore() *PatternStore {
	return &PatternStore{
		patterns: patterns,
	}
}

var patterns = []*api.Pattern{
	&api.Pattern{
		ID:          api.PatternID(uuid.Must(uuid.FromString("a78f6e8d-d0f9-46a0-a8d2-e164dda4bd2b"))),
		Type:        api.Resource,
		Value:       "*mt.js*",
		Description: "MaxTraffic",
		Tags:        []*api.Tag{},
		CreatedAt:   nil,
		DeletedAt:   nil,
	},
	&api.Pattern{
		ID:          api.PatternID(uuid.Must(uuid.FromString("e84b395d-0455-4e27-85dd-211accdc2d4e"))),
		Type:        api.Resource,
		Value:       "https://www.google-analytics.com/analytics.js*",
		Description: "Google analytics",
		Tags:        []*api.Tag{},
		CreatedAt:   nil,
		DeletedAt:   nil,
	},
	&api.Pattern{
		ID:          api.PatternID(uuid.Must(uuid.FromString("16d26858-ca9d-4f34-9f22-0f9245748459"))),
		Type:        api.Resource,
		Value:       "*fbevents.js*",
		Description: "FB analytics",
		Tags:        []*api.Tag{},
		CreatedAt:   nil,
		DeletedAt:   nil,
	},
}
