package mysql

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/smilga/analyzer/api"
)

type WebsiteStore struct {
	DB *sqlx.DB
}

func (s *WebsiteStore) ByUser(id api.UserID) ([]*api.Website, error) {
	ws := []*api.Website{}

	err := s.DB.Select(&ws, "SELECT * FROM websites WHERE user_id=? AND deleted_at IS NULL", id)
	if err != nil {
		return nil, err
	}

	return ws, nil
}

func (s *WebsiteStore) Get(id api.WebsiteID) (*api.Website, error) {
	w := &api.Website{}

	err := s.DB.Get(w, "SELECT * FROM websites where id=? AND deleted_at IS NULL", id)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (s *WebsiteStore) Save(w *api.Website) error {
	now := time.Now()
	if w.ID == 0 {
		w.CreatedAt = &now
	}

	res, err := s.DB.Exec(`
		INSERT INTO websites
		(id, user_id, url, searched_at, created_at, deleted_at)
		VALUES (?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE
		user_id=VALUES(user_id), url=VALUES(url)
	`, w.ID, w.UserID, w.URL, w.SearchedAt, w.CreatedAt, w.DeletedAt)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	w.ID = api.WebsiteID(id)

	return nil
}

func (s *WebsiteStore) Delete(id api.WebsiteID) error {
	_, err := s.DB.Exec(`UPDATE websites SET deleted_at=NOW() where id=?`, id)
	return err
}

func NewWebsiteStore(DB *sqlx.DB) *WebsiteStore {
	return &WebsiteStore{DB}
}
