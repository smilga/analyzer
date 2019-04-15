package api

import "time"

type MatchID int64

// NOTE When update DB fields dont forget to update archive table and trigger
type Match struct {
	ID        MatchID    `db:"id"`
	PatternID PatternID  `db:"pattern_id"`
	Pattern   *Pattern   `db:"-"`
	WebsiteID WebsiteID  `db:"website_id"`
	ReportID  ReportID   `db:"report_id"`
	Value     string     `db:"value"`
	CreatedAt *time.Time `db:"created_at"`
}
