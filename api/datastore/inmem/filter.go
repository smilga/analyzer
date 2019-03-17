package inmem

import (
	uuid "github.com/satori/go.uuid"
	"github.com/smilga/analyzer/api"
)

type FilterStore struct {
	filters []*api.Filter
}

func (s *FilterStore) Save(target *api.Filter) error {
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
		ID:          api.FilterID(uuid.Must(uuid.FromString("3c1800f3-d0e5-4d27-9780-993f8aa015ab"))),
		Name:        "Marketing Services",
		Description: "",
		Tags:        []*api.Tag{},
	},
	&api.Filter{
		ID:          api.FilterID(uuid.Must(uuid.FromString("6e0862f1-6235-4ed2-a21c-95b8b7a75b10"))),
		Name:        "Analytics providers",
		Description: "",
		Tags:        []*api.Tag{},
	},
}
