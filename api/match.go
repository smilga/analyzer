package api

import "time"

type MatchID int64

type Match struct {
	ID        MatchID    `db:"id"`
	PatternID PatternID  `db:"pattern_id"`
	Pattern   *Pattern   `db:"-"`
	WebsiteID WebsiteID  `db:"website_id"`
	ReportID  ReportID   `db:"inspection_id"`
	Value     string     `db:"value"`
	CreatedAt *time.Time `db:"created_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
