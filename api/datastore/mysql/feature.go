package mysql

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/smilga/analyzer/api"
)

type FeatureStore struct {
	DB *sqlx.DB
}

func (s *FeatureStore) All() ([]*api.Feature, error) {
	fs := []*api.Feature{}
	err := s.DB.Select(&fs, "SELECT * FROM features WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	return fs, nil
}

func (s *FeatureStore) Get(id api.FeatureID) (*api.Feature, error) {
	t := &api.Feature{}

	err := s.DB.Get(t, "SELECT * FROM features where id=? AND deleted_at IS NULL", id)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *FeatureStore) Save(t *api.Feature) error {
	now := time.Now()
	if t.ID == 0 {
		t.CreatedAt = &now
	}

	res, err := s.DB.Exec(`
		INSERT INTO features (id, value, created_at, deleted_at)
		VALUES (?, ?, ?, ?)
	`, t.ID, t.Value, t.CreatedAt, t.DeletedAt)

	if err != nil {
		return err
	}

	if t.ID == 0 {
		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		t.ID = api.FeatureID(id)
	}

	return nil
}

func NewFeatureStore(DB *sqlx.DB) *FeatureStore {
	return &FeatureStore{DB}
}
