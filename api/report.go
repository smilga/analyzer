package api

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type MatchReport struct {
	*Pattern
	Value string
}

type ServiceReport struct {
	*ServiceIdentity
	Match []*MatchReport
}

type ShortReport struct {
	WebsiteID  uuid.UUID
	UserID     uuid.UUID
	SearchedAt *time.Time
	Services   []*ServiceIdentity
}

type Report struct {
	WebsiteID      uuid.UUID
	WebsiteURL     string
	Duration       float64
	UserID         uuid.UUID
	SearchedAt     *time.Time
	ServiceReports []*ServiceReport
}
