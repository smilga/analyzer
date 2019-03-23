package mysql

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/smilga/analyzer/api"
)

type UserStore struct {
	DB *sqlx.DB
}

func (s *UserStore) Save(u *api.User) error {
	now := time.Now()
	if u.ID == 0 {
		u.CreatedAt = &now
	} else {
		u.UpdatedAt = &now
	}

	res, err := s.DB.Exec(`
		INSERT INTO users
		(id, name, email, password, created_at, updated_at, deleted_at)
		VALUES (?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE
		name=VALUES(name), updated_at=VALUES(updated_at)
	`, u.ID, u.Name, u.Email, u.Password, u.CreatedAt, u.UpdatedAt, u.DeletedAt)

	if err != nil {
		return err
	}

	if u.ID == 0 {
		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		u.ID = api.UserID(id)
	}

	return nil
}

func (s *UserStore) ByID(id api.UserID) (*api.User, error) {
	u := &api.User{}
	err := s.DB.Get(u, "SELECT * FROM users WHERE id=? AND deleted_at IS NULL", id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *UserStore) ByEmail(email string) (*api.User, error) {
	u := &api.User{}

	err := s.DB.Get(u, "SELECT * FROM users where email=? AND deleted_at IS NULL", email)
	if err == sql.ErrNoRows {
		return nil, api.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func NewUserStore(DB *sqlx.DB) *UserStore {
	return &UserStore{DB}
}
