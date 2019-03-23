package mysql

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
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
		(id, user_id, url, inspected_at, created_at, deleted_at)
		VALUES (?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE
		user_id=VALUES(user_id), url=VALUES(url), inspected_at=VALUES(inspected_at)
	`, w.ID, w.UserID, w.URL, w.InspectedAt, w.CreatedAt, w.DeletedAt)

	if err != nil {
		return err
	}

	if w.ID == 0 {
		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		w.ID = api.WebsiteID(id)
	}

	err = s.storeMatches(w.ID, w.Matches)
	if err != nil {
		return err
	}

	return nil
}

func (s *WebsiteStore) Delete(id api.WebsiteID) error {
	_, err := s.DB.Exec(`UPDATE websites SET deleted_at=NOW() where id=?`, id)
	return err
}

func (s *WebsiteStore) storeMatches(id api.WebsiteID, matches []*api.Match) error {
	// _, err := s.DB.Exec(`UPDATE matches SET deleted_at = NOW() WHERE website_id = ?`, id)
	// if err != nil {
	// 	return err
	// }

	if len(matches) == 0 {
		return nil
	}

	spew.Dump(matches)
	var query string
	values := make([]interface{}, len(matches))
	for i, m := range matches {
		query += fmt.Sprintf("( %d, %d, %d, ?, NOW() )", m.PatternID, m.WebsiteID, m.ReportID)
		if i+1 < len(values) {
			query += ", "
		}
		values[i] = m.Value
	}
	spew.Dump(query)

	_, err := s.DB.Exec(`INSERT INTO matches (pattern_id, website_id, report_id, value, created_at) VALUES `+query, values...)
	return err
}

func NewWebsiteStore(DB *sqlx.DB) *WebsiteStore {
	return &WebsiteStore{DB}
}
