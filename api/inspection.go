package api

import "time"

type InspectionID int64

type Inspection struct {
	ID        InspectionID `db:"id"`
	UserID    UserID       `db:"user_id"`
	WebsiteID WebsiteID    `db:"website_id"`
	ResultID  ResultID     `db:"result_id"`
	CreatedAt *time.Time   `db:"created_at"`
}
