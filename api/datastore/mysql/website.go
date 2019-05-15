package mysql

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/smilga/analyzer/api"
)

type WebsiteStore struct {
	DB *gorm.DB
}

func (s *WebsiteStore) Where(id api.UserID, field string, value interface{}) ([]*api.Website, int, error) {
	ws := []*api.Website{}
	var total int

	err := s.DB.Where("user_id", id).Where(field, value).Find(&ws).Error
	if err != nil {
		return nil, total, err
	}

	err = s.DB.Where("user_id", id).Where(field, value).Count(&total).Error
	if err != nil {
		return nil, total, err
	}

	err = s.AddTags(ws)
	if err != nil {
		return nil, total, err
	}

	return ws, total, nil
}

func (s *WebsiteStore) ByUser(id api.UserID, p *api.Pagination) ([]*api.Website, int, error) {
	ws := []*api.Website{}
	var total int

	q := s.DB.Where("user_id = ?", id)

	if p.ShouldSearch() {
		q.Where(fmt.Sprintf("url LIKE ?", "%%%s%%"), p.Search())
	}

	err := q.Limit(p.Limit()).Offset(p.Offset()).Find(&ws).Error
	if err != nil {
		return nil, total, err
	}

	err = s.DB.Model(&api.Website{}).Where("user_id = ?", id).Count(&total).Error
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

func (s *WebsiteStore) ByFilterID(filterIDs []api.FilterID, id api.UserID, p *api.Pagination) ([]*api.Website, int, error) {
	ws := []*api.Website{}
	var total int

	if len(filterIDs) == 0 {
		return ws, total, nil
	}

	tagIDs := []api.TagID{}
	rows, err := s.DB.DB().Query(fmt.Sprintf("SELECT tag_id FROM filter_tags WHERE filter_id IN (%s);"), filterIDs)
	if err != nil {
		return nil, total, err
	}
	defer rows.Close()
	for rows.Next() {
		var id api.TagID
		err := rows.Scan(&id)
		if err != nil {
			return nil, total, err
		}
		tagIDs = append(tagIDs, id)
	}

	patternIDs := []api.PatternID{}
	rows, err = s.DB.DB().Query(fmt.Sprintf("SELECT pattern_id FROM pattern_tags WHERE tag_id IN (%s);", tagIDs))
	if err != nil {
		return nil, total, err
	}
	defer rows.Close()
	for rows.Next() {
		var id api.PatternID
		err := rows.Scan(&id)
		if err != nil {
			return nil, total, err
		}
		patternIDs = append(patternIDs, id)
	}

	rows, err = s.DB.DB().Query(fmt.Sprintf(`
		SELECT w.* FROM websites w
		INNER JOIN matches m ON m.website_id = w.id
		WHERE m.pattern_id IN (%s)
		AND w.user_id = ?
		AND w.url like ?
		GROUP BY w.id
		LIMIT ?
		OFFSET ?
	`, patternIDs), id, p.Search(), p.Limit(), p.Offset())
	if err != nil {
		return nil, total, err
	}
	defer rows.Close()
	for rows.Next() {
		var w api.Website
		err := rows.Scan(&w.ID, &w.UserID, &w.URL, &w.InspectedAt, &w.CreatedAt, &w.DeletedAt)
		if err != nil {
			return nil, total, err
		}
		ws = append(ws, &w)
	}

	// Add searh dynamicly because of count slow down with "like query"
	countQ := fmt.Sprintf(`
		SELECT COUNT(DISTINCT w.id)
		FROM websites w
		INNER JOIN matches m ON m.website_id = w.id
		WHERE m.pattern_id IN (%s)
		AND w.user_id = ?`, patternIDs)
	countArgs := []interface{}{id}
	if p.ShouldSearch() {
		countQ = fmt.Sprintf("%s AND url like ?", countQ)
		countArgs = append(countArgs, p.Search())
	}

	err = s.DB.DB().QueryRow(countQ, countArgs...).Scan(&total)
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

	err := s.DB.First(&w, int(id)).Error
	if err != nil {
		return nil, err
	}

	// ws := []*api.Website{w}
	// err = s.AddTags(ws)
	// if err != nil {
	// 	return nil, err
	// }

	return w, nil
}

func (s *WebsiteStore) Save(w *api.Website) error {
	now := time.Now()
	if w.ID == 0 {
		w.CreatedAt = &now
	}

	res, err := s.DB.DB().Exec(`
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

	_, err := s.DB.DB().Exec(q, params...)
	if err != nil {
		return err
	}

	return nil
}

func (s *WebsiteStore) Delete(id api.WebsiteID) error {
	// TODO validate id website belongs to user
	// pass context down with user id and then check
	_, err := s.DB.DB().Exec(`UPDATE websites SET deleted_at=NOW() where id=?`, id)
	return err
}

func (s *WebsiteStore) storeMatches(id api.WebsiteID, matches []*api.Match) error {
	// NOTE there is trigger that moves those matches to matches_archive table
	_, err := s.DB.DB().Exec(`DELETE FROM matches WHERE website_id = ?`, id)
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

	_, err = s.DB.DB().Exec(`INSERT INTO matches (pattern_id, website_id, report_id, value, created_at) VALUES `+query, values...)
	return err
}

func (s *WebsiteStore) AddTags(websites []*api.Website) error {
	tags := make(map[api.WebsiteID][]*api.Tag, len(websites))
	ids := make([]string, len(websites))

	for i, w := range websites {
		tags[w.ID] = []*api.Tag{}
		ids[i] = w.ID.String()
	}

	if len(websites) == 0 {
		return nil
	}

	rows, err := s.DB.DB().Query(fmt.Sprintf(`
		SELECT w.id, t.*
		FROM websites w
		INNER JOIN matches m on m.website_id = w.id
		INNER JOIN pattern_tags pt on pt.pattern_id = m.pattern_id
		INNER JOIN tags t on t.id = pt.tag_id
		WHERE w.id in (%s)
	`, strings.Join(ids, ",")))
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

func NewWebsiteStore(DB *gorm.DB) *WebsiteStore {
	return &WebsiteStore{DB}
}
