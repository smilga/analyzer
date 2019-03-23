package api

import "time"

type Result struct {
	Matches         []*Match `json:"matches"`
	LoadedIn        string   `json:"loadedIn"`
	ResourceCheckIn string   `json:"resourceCheckIn"`
	HTMLCheckIn     string   `json:"htmlCheckIn"`
	TotalIn         string   `json:"totalIn"`
	CreatedAt       *time.Time
}
