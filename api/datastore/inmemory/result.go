package inmemory

import (
	uuid "github.com/satori/go.uuid"
	"github.com/smilga/analyzer/api"
)

type ResultStore struct {
	results []*api.Result
}

func (s *ResultStore) ByWebsite(id uuid.UUID) ([]*api.Result, error) {
	results := []*api.Result{}
	for _, r := range s.results {
		if r.WebsiteID == id {
			results = append(results, r)
		}
	}
	return results, nil
}

func (s *ResultStore) LatestByWebsite(id uuid.UUID) (*api.Result, error) {
	result := &api.Result{}
	for _, r := range s.results {
		if r.WebsiteID == id && (result.CreatedAt == nil || result.CreatedAt.Before(*r.CreatedAt)) {
			result = r
		}
	}
	return result, nil
}

func (s *ResultStore) Save(r *api.Result) error {
	s.results = append(s.results, r)
	return nil
}

func NewResultStore() *ResultStore {
	return &ResultStore{}
}
