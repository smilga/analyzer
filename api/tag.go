package api

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Tag struct {
	Name      string
	UserID    uuid.UUID
	CreatedAt *time.Time
}
