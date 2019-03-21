package mysql

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/smilga/analyzer/api"
)

type PatternStore struct {
	DB *sqlx.DB
}

func (s *PatternStore) Save(p *api.Pattern) error {
	now := time.Now()
	if p.ID == 0 {
		p.CreatedAt = &now
	} else {
		p.UpdatedAt = &now
	}

	res, err := s.DB.Exec(`
		INSERT INTO patterns
		(id, type, value, description, created_at, updated_at, deleted_at)
		VALUES (?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE
		type=VALUES(type), description=VAlUES(description), updated_at=VALUES(updated_at)
	`, p.ID, p.Type, p.Value, p.Description, p.CreatedAt, p.UpdatedAt, p.DeletedAt)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = api.PatternID(id)

	if len(p.Tags) > 0 {
		err = s.storeTags(p.ID, p.Tags)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PatternStore) Get(id api.PatternID) (*api.Pattern, error) {
	p := &api.Pattern{}
	rows, err := s.DB.Query(`
		SELECT p.*, t.* FROM patterns p
		LEFT JOIN pattern_tags pt on pt.pattern_id = p.id
		LEFT JOIN tags t on pt.tag_id = t.id
		WHERE p.id = ?
		AND p.deleted_at IS NULL
	`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		t := &api.Tag{}
		var nullID *int64
		var nullValue *string
		err := rows.Scan(&p.ID, &p.Type, &p.Value, &p.Description, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt, &nullID, &nullValue, &t.CreatedAt, &t.DeletedAt)
		if err != nil {
			return nil, err
		}
		if nullID != nil {
			t.ID = api.TagID(*nullID)
			t.Value = *nullValue
			p.Tags = append(p.Tags, t)
		}
	}

	return p, nil
}

func (s *PatternStore) All() ([]*api.Pattern, error) {
	ps := []*api.Pattern{}
	rows, err := s.DB.Query(`
		SELECT p.*, t.* FROM patterns p
		LEFT JOIN pattern_tags pt on pt.pattern_id = p.id
		LEFT JOIN tags t on pt.tag_id = t.id
		WHERE p.deleted_at IS NULL
	`)
	if err != nil {
		return nil, err
	}

	patterns := make(map[api.PatternID]*api.Pattern)
	patternTags := make(map[api.PatternID][]*api.Tag)
	for rows.Next() {
		p := &api.Pattern{}
		t := &api.Tag{}
		var nullID *int64
		var nullValue *string
		err := rows.Scan(&p.ID, &p.Type, &p.Value, &p.Description, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt, &nullID, &nullValue, &t.CreatedAt, &t.DeletedAt)
		if err != nil {
			return nil, err
		}
		patterns[p.ID] = p
		if nullID != nil {
			t.ID = api.TagID(*nullID)
			t.Value = *nullValue
			patternTags[p.ID] = append(patternTags[p.ID], t)
		}
	}

	for _, p := range patterns {
		if tags, ok := patternTags[p.ID]; ok {
			p.Tags = tags
		} else {
			p.Tags = []*api.Tag{}
		}
		ps = append(ps, p)
	}

	return ps, nil
}

func (s *PatternStore) Delete(id api.PatternID) error {
	_, err := s.DB.Exec(`UPDATE patterns SET deleted_at=NOW() where id=?`, id)
	return err
}

func (s *PatternStore) storeTags(id api.PatternID, tags []*api.Tag) error {
	var bindVars string
	tagIDs := make([]interface{}, len(tags))
	for i, t := range tags {
		bindVars += fmt.Sprintf("( %d, ? )", id)
		if i+1 < len(tagIDs) {
			bindVars += ", "
		}
		tagIDs[i] = t.ID
	}

	_, err := s.DB.Exec(`DELETE FROM pattern_tags WHERE pattern_id = ?`, id)
	if err != nil {
		return err
	}
	_, err = s.DB.Exec(`INSERT INTO pattern_tags (pattern_id, tag_id) VALUES `+bindVars, tagIDs...)
	return err
}

func NewPatternStore(DB *sqlx.DB) *PatternStore {
	return &PatternStore{DB}
}
