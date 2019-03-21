package api

import "time"

type ResultID int64

type ResultStorage interface {
	Save(*Result) error
	LatestByWebsite(WebsiteID) (*Result, error)
}

type Result struct {
	ID ResultID `db:"id"`
	*ResponseTime
	Matches   []*Match   `db:"-"`
	CreatedAt *time.Time `db:"created_at"`
}

type ResponseTime struct {
	LoadedIn        string `db:"loaded_in"`
	ResourceCheckIn string `db:"resource_check_in"`
	HTMLCheckIn     string `db:"html_check_in"`
	TotalIn         string `db:"total_in"`
}
