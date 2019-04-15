package api

import "time"

type ServiceID int64

type Service struct {
	ID        ServiceID  `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Features  []*Feature `db:"-" json:"features"`
	CreatedAt *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"`
}

type ServiceStorage interface {
	Save(*Service) error
	Delete(ServiceID) error
	Get(ServiceID) (*Service, error)
	All(*Pagination) ([]*Service, int, error)
	ByFeatures([]FeatureID, *Pagination) ([]*Service, int, error)
}
