package api

import (
	"fmt"
	"time"
)

func WebsitesToCsv(websites []*Website) [][]string {
	cols := make([][]string, len(websites)+1)

	cols[0] = []string{"Website", "InspectedAt", "Tags"}
	for i, w := range websites {
		cols[i+1] = []string{w.URL, timeFormat(w.InspectedAt), tagsString(w.Tags)}
	}

	return cols
}

func tagsString(tags []*Tag) string {
	if tags == nil {
		return ""
	}
	var ts string
	for _, t := range tags {
		ts += fmt.Sprintf("%s ", t.Value)
	}
	return ts
}

func timeFormat(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
