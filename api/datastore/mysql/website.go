package mysql

import (
	"fmt"
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

	err = s.AddTags(ws)
	if err != nil {
		return nil, err
	}

	return ws, nil
}

func (s *WebsiteStore) ByFilterID(filterIDs []api.FilterID, id api.UserID) ([]*api.Website, error) {
	ws := []*api.Website{}

	if len(filterIDs) == 0 {
		return ws, nil
	}

	tagIDs := []api.TagID{}
	query, args, err := sqlx.In("SELECT tag_id FROM filter_tags WHERE filter_id IN (?);", filterIDs)
	if err != nil {
		return nil, err
	}
	err = s.DB.Select(&tagIDs, query, args...)
	if err != nil {
		return nil, err
	}

	patternIDs := []api.PatternID{}
	query, args, err = sqlx.In("SELECT pattern_id FROM pattern_tags WHERE tag_id IN (?);", tagIDs)
	if err != nil {
		return nil, err
	}
	err = s.DB.Select(&patternIDs, query, args...)
	if err != nil {
		return nil, err
	}

	query, args, err = sqlx.In(`
		SELECT w.* from matches m
		LEFT JOIN websites w
		ON m.website_id = w.id
		WHERE m.pattern_id IN (?)
		AND m.deleted_at IS NULL
		GROUP BY w.id
		HAVING w.user_id = ?;
	`, patternIDs, id)
	if err != nil {
		return nil, err
	}
	err = s.DB.Select(&ws, query, args...)
	if err != nil {
		return nil, err
	}

	err = s.AddTags(ws)
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

	ws := []*api.Website{w}
	err = s.AddTags(ws)
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
	_, err := s.DB.Exec(`UPDATE matches SET deleted_at = NOW() WHERE website_id = ?`, id)
	if err != nil {
		return err
	}

	if len(matches) == 0 {
		return nil
	}

	var query string
	values := make([]interface{}, len(matches))
	for i, m := range matches {
		query += fmt.Sprintf("( %d, %d, %d, ?, NOW() )", m.PatternID, m.WebsiteID, m.ReportID)
		if i+1 < len(values) {
			query += ", "
		}
		values[i] = m.Value
	}

	_, err = s.DB.Exec(`INSERT INTO matches (pattern_id, website_id, report_id, value, created_at) VALUES `+query, values...)
	return err
}

func (s *WebsiteStore) AddTags(websites []*api.Website) error {
	tags := make(map[api.WebsiteID][]*api.Tag, len(websites))
	websiteIDs := make([]api.WebsiteID, len(websites))

	for i, w := range websites {
		tags[w.ID] = []*api.Tag{}
		websiteIDs[i] = w.ID
	}

	if len(websites) == 0 {
		return nil
	}

	query, args, err := sqlx.In(`
		SELECT w.id, t.* FROM websites w
		INNER JOIN matches m on m.website_id = w.id
		INNER JOIN pattern_tags pt on pt.pattern_id = m.pattern_id
		INNER JOIN tags t on t.id = pt.tag_id WHERE w.id in (?)
		AND m.deleted_at IS NULL;
	`, websiteIDs)
	if err != nil {
		return err
	}

	rows, err := s.DB.Query(query, args...)
	if err != nil {
		return err
	}

	for rows.Next() {
		var wID api.WebsiteID
		t := &api.Tag{}
		err := rows.Scan(&wID, &t.ID, &t.Value, &t.CreatedAt, &t.DeletedAt)
		if err != nil {
			return err
		}
		tags[wID] = append(tags[wID], t)
	}

	for _, w := range websites {
		if ts, ok := tags[w.ID]; ok {
			w.Tags = ts
		}
	}
	return nil
}

func NewWebsiteStore(DB *sqlx.DB) *WebsiteStore {
	return &WebsiteStore{DB}
}
