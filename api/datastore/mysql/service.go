package mysql

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	"github.com/smilga/analyzer/api"
)

type ServiceStore struct {
	DB *gorm.DB
}

func (s *ServiceStore) Save(p *api.Service) error {
	err := s.DB.Create(p).Error
	if err != nil {
		return err
	}

	err = s.updateServiceFeatures(p.ID, p.Features)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceStore) Get(id api.ServiceID) (*api.Service, error) {
	p := &api.Service{}
	rows, err := s.DB.DB().Query(`
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
	ss := []*api.Service{}
	var total int

	err := s.DB.Limit(p.Limit()).Offset(p.Offset()).Find(&ss).Error
	if err != nil {
		return nil, total, err
	}

	err = s.DB.Model(&api.Service{}).Count(&total).Error
	if err != nil {
		return nil, total, err
	}

	if !p.NoLimit() {
		err := s.AddFeatures(ss)
		if err != nil {
			return nil, total, err
		}
	}

	return ss, total, nil
}

func (s *ServiceStore) ByFeatures(featureIDs []api.FeatureID, p *api.Pagination) ([]*api.Service, int, error) {
	ws := []*api.Service{}
	var total int

	rows, err := s.DB.DB().Query(fmt.Sprintf(`
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
	`, featureIDs), p.Search(), p.Limit(), p.Offset())
	if err != nil {
		return nil, total, err
	}

	defer rows.Close()

	for rows.Next() {
		var s api.Service
		err := rows.Scan(&s.ID, &s.Name, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt)
		if err != nil {
			return nil, total, err
		}
		ws = append(ws, &s)
	}

	err = s.DB.DB().QueryRow(fmt.Sprintf(`
		SELECT count(*) FROM services
		WHERE id IN (
			SELECT service_id FROM service_features
			WHERE feature_id in (%s)
			GROUP BY id
		)
		AND deleted_at IS NULL
		AND name like ?
	`, featureIDs), p.Search()).Scan(&total)
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
	_, err = s.DB.DB().Exec(`UPDATE services SET deleted_at=NOW() where id=?`, id)
	return err
}

func (s *ServiceStore) updateServiceFeatures(id api.ServiceID, features []*api.Feature) error {
	_, err := s.DB.DB().Exec(`DELETE FROM service_features WHERE service_id = ?`, id)
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

	_, err = s.DB.DB().Exec(`INSERT INTO service_features (service_id, feature_id) VALUES `+bindVars, featureIDs...)
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

	rows, err := s.DB.DB().Query(query, args...)
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

func NewServiceStore(DB *gorm.DB) *ServiceStore {
	return &ServiceStore{DB}
}
