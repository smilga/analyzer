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

func (s *WebsiteStore) ByUser(id api.UserID, p *api.Pagination) ([]*api.Website, int, error) {
	ws := []*api.Website{}
	var total int

	err := s.DB.Select(&ws, `
		SELECT * FROM websites
		WHERE user_id=?
		AND deleted_at IS NULL
		AND url like ?
		LIMIT ?
		OFFSET ?
		`, id, p.Search(), p.Limit(), p.Offset())
	if err != nil {
		return nil, total, err
	}

	err = s.DB.Get(&total, `
		SELECT count(*) FROM websites
		WHERE user_id=?
		AND deleted_at IS NULL
		AND url like ?
		`, id, p.Search())
	if err != nil {
		return nil, total, err
	}

	if !p.NoLimit() {
		err = s.AddTags(ws)
		if err != nil {
			return nil, total, err
		}
	}

	return ws, total, nil
}

func (s *WebsiteStore) Where(id api.UserID, field string, value interface{}) ([]*api.Website, int, error) {
	ws := []*api.Website{}
	var total int

	var query string
	if value == "NULL" {
		query = fmt.Sprintf("%s IS NULL", field)
	} else {
		query = fmt.Sprintf("%s = %v", field, value)
	}

	err := s.DB.Select(&ws, fmt.Sprintf(`
		SELECT * FROM websites
		WHERE user_id=?
		AND deleted_at IS NULL
		AND %s
		`, query), id)
	if err != nil {
		return nil, total, err
	}

	err = s.DB.Get(&total, fmt.Sprintf(`
		SELECT count(*) FROM websites
		WHERE user_id=?
		AND deleted_at IS NULL
		AND %s
		`, query), id)
	if err != nil {
		return nil, total, err
	}

	return ws, total, nil
}

func (s *WebsiteStore) ByFilterID(filterIDs []api.FilterID, id api.UserID, p *api.Pagination) ([]*api.Website, int, error) {
	ws := []*api.Website{}
	var total int

	if len(filterIDs) == 0 {
		return ws, total, nil
	}

	tagIDs := []api.TagID{}
	query, args, err := sqlx.In("SELECT tag_id FROM filter_tags WHERE filter_id IN (?);", filterIDs)
	if err != nil {
		return nil, total, err
	}
	err = s.DB.Select(&tagIDs, query, args...)
	if err != nil {
		return nil, total, err
	}

	patternIDs := []api.PatternID{}
	query, args, err = sqlx.In("SELECT pattern_id FROM pattern_tags WHERE tag_id IN (?);", tagIDs)
	if err != nil {
		return nil, total, err
	}
	err = s.DB.Select(&patternIDs, query, args...)
	if err != nil {
		return nil, total, err
	}
	// SELECT w.* from websites w WHERE w.id IN (SELECT website_id from matches where pattern_id in (?) GROUP BY website_id) AND w.user_id = ? AND w.url like '%%' LIMIT 10 OFFSET 0;
	query, args, err = sqlx.In(`
		SELECT w.* from websites w
		RIGHT JOIN matches m ON m.website_id = w.id
		WHERE m.pattern_id IN (?)
		AND w.user_id = ?
		AND w.url like ?
		LIMIT ?
		OFFSET ?
	`, patternIDs, id, p.Search(), p.Limit(), p.Offset())
	if err != nil {
		return nil, total, err
	}
	err = s.DB.Select(&ws, query, args...)
	if err != nil {
		return nil, total, err
	}
	// SELECT count(*) from websites w where w.id in (SELECT website_id from matches where pattern_id in (1,2) GROUP BY website_id) AND w.user_id = ? AND w.url like '%%';
	// SELECT count(DISTINCT w.id) FROM websites w INNER JOIN matches m ON m.website_id = w.id WHERE m.pattern_id IN (?) AND w.user_id = ?;
	// add like search dinamically
	query, args, err = sqlx.In(`
		SELECT count(*)
		FROM (SELECT count(*) FROM websites w
		RIGHT JOIN matches m ON m.website_id = w.id
		WHERE m.pattern_id IN (?)
		AND w.user_id = ?
		AND w.url like ?
		GROUP BY w.id) as total`, patternIDs, id, p.Search())
	if err != nil {
		return nil, total, err
	}
	err = s.DB.Get(&total, query, args...)
	if err != nil {
		return nil, total, err
	}

	if !p.NoLimit() {
		err = s.AddTags(ws)
		if err != nil {
			return nil, total, err
		}
	}

	return ws, total, nil
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

func (s *WebsiteStore) SaveBatch(websites []*api.Website) error {
	batchLimit := 20000

	if len(websites) > batchLimit {
		chunks := int(len(websites) / batchLimit)
		left := len(websites) % batchLimit

		for i := 0; i < chunks; i++ {
			err := s.saveBatch(websites[i*batchLimit : (i+1)*batchLimit])
			if err != nil {
				return err
			}
		}
		err := s.saveBatch(websites[chunks*batchLimit : chunks*batchLimit+left])
		if err != nil {
			return err
		}

		return nil
	}

	return s.saveBatch(websites)
}

func (s *WebsiteStore) saveBatch(websites []*api.Website) error {
	var params []interface{}
	q := "INSERT INTO websites (user_id, url, created_at) VALUES"
	for i, w := range websites {
		q += "(?,?,?)"
		if i+1 < len(websites) {
			q += ","
		}
		params = append(params, w.UserID, w.URL, w.CreatedAt)
	}

	_, err := s.DB.Exec(q, params...)
	if err != nil {
		return err
	}

	return nil
}

func (s *WebsiteStore) Delete(id api.WebsiteID) error {
	// TODO validate id website belongs to user
	// pass context down with user id and then check
	_, err := s.DB.Exec(`UPDATE websites SET deleted_at=NOW() where id=?`, id)
	return err
}

func (s *WebsiteStore) storeMatches(id api.WebsiteID, matches []*api.Match) error {
	_, err := s.DB.Exec(`DELETE from matches WHERE website_id = ?`, id)
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
		SELECT w.id, t.* FROM websites w INNER JOIN matches m on m.website_id = w.id INNER JOIN pattern_tags pt on pt.pattern_id = m.pattern_id INNER JOIN tags t on t.id = pt.tag_id WHERE w.id in (?) AND m.deleted_at IS NULL;
	`, websiteIDs)
	if err != nil {
		return err
	}

	rows, err := s.DB.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

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
