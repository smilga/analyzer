package mysql

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/smilga/analyzer/api"
)

type FilterStore struct {
	DB *sqlx.DB
}

func (s *FilterStore) Get(id api.FilterID) (*api.Filter, error) {
	f := &api.Filter{}
	rows, err := s.DB.Query(`
		SELECT f.*, t.* FROM filter f
		LEFT JOIN filter_tags ft on ft.filter_id = f.id
		LEFT JOIN tags t on ft.tag_id = t.id
		WHERE f.id = ?
		AND f.deleted_at IS NULL
	`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		t := &api.Tag{}
		var nullID *int64
		var nullValue *string
		err := rows.Scan(&f.ID, &f.Name, &f.Description, &f.CreatedAt, &f.UpdatedAt, &f.DeletedAt, &nullID, &nullValue, &t.CreatedAt, &t.DeletedAt)
		if err != nil {
			return nil, err
		}
		if nullID != nil {
			t.ID = api.TagID(*nullID)
			t.Value = *nullValue
			f.Tags = append(f.Tags, t)
		}
	}

	return f, nil
}

func (s *FilterStore) All() ([]*api.Filter, error) {
	fs := []*api.Filter{}
	rows, err := s.DB.Query(`
		SELECT f.*, t.* FROM filters f
		LEFT JOIN filter_tags ft on ft.filter_id = f.id
		LEFT JOIN tags t on ft.tag_id = t.id
		WHERE f.deleted_at IS NULL
	`)
	if err != nil {
		return nil, err
	}

	filters := make(map[api.FilterID]*api.Filter)
	filterTags := make(map[api.FilterID][]*api.Tag)
	for rows.Next() {
		f := &api.Filter{}
		t := &api.Tag{}
		var nullID *int64
		var nullValue *string
		err := rows.Scan(&f.ID, &f.Name, &f.Description, &f.CreatedAt, &f.UpdatedAt, &f.DeletedAt, &nullID, &nullValue, &t.CreatedAt, &t.DeletedAt)
		if err != nil {
			return nil, err
		}
		filters[f.ID] = f
		if nullID != nil {
			t.ID = api.TagID(*nullID)
			t.Value = *nullValue
			filterTags[f.ID] = append(filterTags[f.ID], t)
		}
	}

	for _, f := range filters {
		if tags, ok := filterTags[f.ID]; ok {
			f.Tags = tags
		} else {
			f.Tags = []*api.Tag{}
		}
		fs = append(fs, f)
	}

	return fs, nil
}

func (s *FilterStore) Save(f *api.Filter) error {
	now := time.Now()
	if f.ID == 0 {
		f.CreatedAt = &now
	} else {
		f.UpdatedAt = &now
	}

	res, err := s.DB.Exec(`
		INSERT INTO filters
		(id, name, description, created_at, updated_at, deleted_at)
		VALUES (?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE
		name=VALUES(name), description=VAlUES(description), updated_at=VALUES(updated_at)
	`, f.ID, f.Name, f.Description, f.CreatedAt, f.UpdatedAt, f.DeletedAt)
	if err != nil {
		return err
	}

	if f.ID == 0 {
		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		f.ID = api.FilterID(id)
	}

	err = s.updateFilterTags(f.ID, f.Tags)
	if err != nil {
		return err
	}

	return nil
}

func (s *FilterStore) Delete(id api.FilterID) error {
	_, err := s.DB.Exec(`UPDATE filters SET deleted_at=NOW() where id=?`, id)
	return err
}

func (s *FilterStore) updateFilterTags(id api.FilterID, tags []*api.Tag) error {
	_, err := s.DB.Exec(`DELETE FROM filter_tags WHERE filter_id = ?`, id)
	if err != nil {
		return err
	}

	if len(tags) == 0 {
		return nil
	}

	var bindVars string
	tagIDs := make([]interface{}, len(tags))
	for i, t := range tags {
		bindVars += fmt.Sprintf("( %d, ? )", id)
		if i+1 < len(tagIDs) {
			bindVars += ", "
		}
		tagIDs[i] = t.ID
	}

	_, err = s.DB.Exec(`INSERT INTO filter_tags (filter_id, tag_id) VALUES `+bindVars, tagIDs...)
	return err
}

func NewFilterStore(DB *sqlx.DB) *FilterStore {
	return &FilterStore{DB}
}
