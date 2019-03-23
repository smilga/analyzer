package mysql

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/smilga/analyzer/api"
)

type TagStore struct {
	DB *sqlx.DB
}

func (s *TagStore) All() ([]*api.Tag, error) {
	ts := []*api.Tag{}
	err := s.DB.Select(&ts, "SELECT * FROM tags WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func (s *TagStore) Get(id api.TagID) (*api.Tag, error) {
	t := &api.Tag{}

	err := s.DB.Get(t, "SELECT * FROM tags where id=? AND deleted_at IS NULL", id)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TagStore) Save(t *api.Tag) error {
	now := time.Now()
	if t.ID == 0 {
		t.CreatedAt = &now
	}

	res, err := s.DB.Exec(`
		INSERT INTO tags (id, value, created_at, deleted_at)
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
		t.ID = api.TagID(id)
	}

	return nil
}

func NewTagStore(DB *sqlx.DB) *TagStore {
	return &TagStore{DB}
}
