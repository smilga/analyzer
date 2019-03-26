package api

import "time"

type Result struct {
	WebsiteID       WebsiteID `json:"websiteId"`
	UserID          UserID    `json:"userId"`
	Matches         []*Match  `json:"matches"`
	StartedIn       string    `json:"startedIn"`
	LoadedIn        string    `json:"loadedIn"`
	ResourceCheckIn string    `json:"resourceCheckIn"`
	HTMLCheckIn     string    `json:"htmlCheckIn"`
	TotalIn         string    `json:"totalIn"`
	CreatedAt       *time.Time
}
