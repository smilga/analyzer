package api

import "time"

type ReportID int64

type ReportStorage interface {
	Save(*Report) error
	ByWebsite(WebsiteID) (*Report, error)
}

type Report struct {
	ID              ReportID   `db:"id"`
	UserID          UserID     `db:"user_id"`
	WebsiteID       WebsiteID  `db:"website_id"`
	Matches         []*Match   `db:"-"`
	WebsiteURL      string     `db:"-"`
	LoadedIn        string     `db:"loaded_in"`
	StartedIn       string     `db:"started_in"`
	ResourceCheckIn string     `db:"resource_check_in"`
	HTMLCheckIn     string     `db:"html_check_in"`
	TotalIn         string     `db:"total_in"`
	CreatedAt       *time.Time `db:"created_at"`
}
