package mysql

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/smilga/analyzer/api"
)

type ServiceStore struct {
	DB *sqlx.DB
}

func (s *ServiceStore) Save(p *api.Service) error {
	now := time.Now()
	if p.ID == 0 {
		p.CreatedAt = &now
	} else {
		p.UpdatedAt = &now
	}

	res, err := s.DB.Exec(`
		INSERT INTO services
		(id, name, created_at, updated_at, deleted_at)
		VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE
		name=VALUES(name), updated_at=VALUES(updated_at)
	`, p.ID, p.Name, p.CreatedAt, p.UpdatedAt, p.DeletedAt)
	if err != nil {
		return err
	}

	if p.ID == 0 {
		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		p.ID = api.ServiceID(id)
	}

	err = s.updateServiceFeatures(p.ID, p.Features)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceStore) Get(id api.ServiceID) (*api.Service, error) {
	p := &api.Service{}
	rows, err := s.DB.Query(`
		SELECT s.*, f.* FROM services s
		LEFT JOIN service_features sf on sf.service_id = s.id
		LEFT JOIN features f on sf.feature_id = f.id
		WHERE s.id = ?
		AND s.deleted_at IS NULL
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		t := &api.Feature{}
		var nullID *int64
		var nullValue *string
		err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt, &nullID, &nullValue, &t.CreatedAt, &t.DeletedAt)
		if err != nil {
			return nil, err
		}
		if nullID != nil {
			t.ID = api.FeatureID(*nullID)
			t.Value = *nullValue
			p.Features = append(p.Features, t)
		}
	}

	return p, nil
}

func (s *ServiceStore) All(p *api.Pagination) ([]*api.Service, int, error) {
	ws := []*api.Service{}
	var total int

	err := s.DB.Select(&ws, `
		SELECT * FROM services
		WHERE deleted_at IS NULL
		AND name like ?
		LIMIT ?
		OFFSET ?
		`, p.Search(), p.Limit(), p.Offset())
	if err != nil {
		return nil, total, err
	}

	err = s.DB.Get(&total, `
		SELECT count(*) FROM services
		WHERE deleted_at IS NULL
		AND name like ?
		`, p.Search())
	if err != nil {
		return nil, total, err
	}

	if !p.NoLimit() {
		err := s.AddFeatures(ws)
		if err != nil {
			return nil, total, err
		}
	}

	return ws, total, nil
}

func (s *ServiceStore) ByFeatures(featureIDs []api.FeatureID, p *api.Pagination) ([]*api.Service, int, error) {
	ws := []*api.Service{}
	var total int

	query, args, err := sqlx.In(`
		SELECT * FROM services
		WHERE id IN (
			SELECT service_id FROM service_features
			WHERE feature_id in (?)
			GROUP BY id
		)
		AND deleted_at IS NULL
		AND name like ?
		LIMIT ?
		OFFSET ?
	`, featureIDs, p.Search(), p.Limit(), p.Offset())
	if err != nil {
		return nil, total, err
	}
	err = s.DB.Select(&ws, query, args...)
	if err != nil {
		return nil, total, err
	}

	query, args, err = sqlx.In(`
		SELECT count(*) FROM services
		WHERE id IN (
			SELECT service_id FROM service_features
			WHERE feature_id in (?)
			GROUP BY id
		)
		AND deleted_at IS NULL
		AND name like ?
	`, featureIDs, p.Search())
	if err != nil {
		return nil, total, err
	}
	err = s.DB.Get(&total, query, args...)
	if err != nil {
		return nil, total, err
	}

	if !p.NoLimit() {
		err := s.AddFeatures(ws)
		if err != nil {
			return nil, total, err
		}
	}

	return ws, total, nil
}

func (s *ServiceStore) Delete(id api.ServiceID) error {
	t := []*api.Feature{}
	err := s.updateServiceFeatures(id, t)
	if err != nil {
		return err
	}
	_, err = s.DB.Exec(`UPDATE services SET deleted_at=NOW() where id=?`, id)
	return err
}

func (s *ServiceStore) updateServiceFeatures(id api.ServiceID, features []*api.Feature) error {
	_, err := s.DB.Exec(`DELETE FROM service_features WHERE service_id = ?`, id)
	if err != nil {
		return err
	}

	if len(features) == 0 {
		return nil
	}

	var bindVars string
	featureIDs := make([]interface{}, len(features))
	for i, t := range features {
		bindVars += fmt.Sprintf("( %d, ? )", id)
		if i+1 < len(featureIDs) {
			bindVars += ", "
		}
		featureIDs[i] = t.ID
	}

	_, err = s.DB.Exec(`INSERT INTO service_features (service_id, feature_id) VALUES `+bindVars, featureIDs...)
	return err
}

func (s *ServiceStore) AddFeatures(services []*api.Service) error {
	features := make(map[api.ServiceID][]*api.Feature, len(services))
	serviceIDs := make([]api.ServiceID, len(services))

	for i, s := range services {
		features[s.ID] = []*api.Feature{}
		serviceIDs[i] = s.ID
	}

	if len(services) == 0 {
		return nil
	}

	query, args, err := sqlx.In(`
		SELECT s.id, f.*
		FROM services s
		JOIN service_features sf on sf.service_id = s.id
		JOIN features f on f.id = sf.feature_id
		WHERE s.id in (?)
	`, serviceIDs)
	if err != nil {
		return err
	}

	rows, err := s.DB.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var sID api.ServiceID
		f := &api.Feature{}
		err := rows.Scan(&sID, &f.ID, &f.Value, &f.CreatedAt, &f.DeletedAt)
		if err != nil {
			return err
		}
		features[sID] = append(features[sID], f)
	}

	for _, s := range services {
		if ts, ok := features[s.ID]; ok {
			s.Features = ts
		}
	}
	return nil
}

func NewServiceStore(DB *sqlx.DB) *ServiceStore {
	return &ServiceStore{DB}
}
