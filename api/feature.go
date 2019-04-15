package api

import (
	"errors"
	"time"
)

var (
	ErrFeatureNotFound = errors.New("Errot feature not found")
)

type FeatureID int64

type FeatureStorage interface {
	Get(FeatureID) (*Feature, error)
	All() ([]*Feature, error)
	Save(*Feature) error
}

type Feature struct {
	ID        FeatureID  `db:"id" json:"id"`
	Value     string     `db:"value" json:"value"`
	CreatedAt *time.Time `db:"created_at" json:"createdAt"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"`
}

func NewFeature(value string) *Feature {
	now := time.Now()

	return &Feature{
		Value:     value,
		CreatedAt: &now,
		DeletedAt: nil,
	}
}
