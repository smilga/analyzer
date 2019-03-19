package inmem

import (
	"github.com/smilga/analyzer/api"
)

type FilterStore struct {
	filters []*api.Filter
}

func (s *FilterStore) Save(target *api.Filter) error {
	if target.ID == 0 {
		var last int64
		for _, n := range s.filters {
			if int64(n.ID) > last {
				last = int64(n.ID)
			}
		}
		target.ID = api.FilterID(last + 1)
	}

	for i, t := range s.filters {
		if t.ID == api.FilterID(target.ID) {
			s.filters = append(s.filters[:i], s.filters[i+1:]...)
		}
	}
	s.filters = append(s.filters, target)

	return nil
}

func (s *FilterStore) All() ([]*api.Filter, error) {
	return s.filters, nil
}

func (s *FilterStore) Get(id api.FilterID) (*api.Filter, error) {
	for _, t := range s.filters {
		if t.ID == api.FilterID(id) {
			return t, nil
		}
	}
	return nil, api.ErrFilterNotFound
}

func NewFilterStore() *FilterStore {
	return &FilterStore{
		filters: filters,
	}
}

var filters = []*api.Filter{
	&api.Filter{
		ID:          1,
		Name:        "Marketing Services",
		Description: "",
		Tags:        []*api.Tag{},
	},
	&api.Filter{
		ID:          2,
		Name:        "Analytics providers",
		Description: "",
		Tags:        []*api.Tag{},
	},
}
