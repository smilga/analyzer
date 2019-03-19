package inmem

import (
	"github.com/smilga/analyzer/api"
)

type PatternStore struct {
	patterns []*api.Pattern
}

func (s *PatternStore) Save(target *api.Pattern) error {
	if target.ID == 0 {
		var last int64
		for _, n := range s.patterns {
			if int64(n.ID) > last {
				last = int64(n.ID)
			}
		}
		target.ID = api.PatternID(last + 1)
	}

	for i, p := range s.patterns {
		if p.ID == target.ID {
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
		ID:          1,
		Type:        api.Resource,
		Value:       "*mt.js*",
		Description: "MaxTraffic",
		Tags:        []*api.Tag{},
		CreatedAt:   nil,
		DeletedAt:   nil,
	},
	&api.Pattern{
		ID:          2,
		Type:        api.Resource,
		Value:       "https://www.google-analytics.com/analytics.js*",
		Description: "Google analytics",
		Tags:        []*api.Tag{},
		CreatedAt:   nil,
		DeletedAt:   nil,
	},
	&api.Pattern{
		ID:          3,
		Type:        api.Resource,
		Value:       "*fbevents.js*",
		Description: "FB analytics",
		Tags:        []*api.Tag{},
		CreatedAt:   nil,
		DeletedAt:   nil,
	},
	&api.Pattern{
		ID:          4,
		Type:        api.HTML,
		Value:       ".top-menu-item",
		Description: "Matches given.lv toolbar item",
		Tags:        []*api.Tag{},
		CreatedAt:   nil,
		DeletedAt:   nil,
	},
}
