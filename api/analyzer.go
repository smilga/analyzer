package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type Analyzer struct {
	PatternStorage PatternStorage
	WebsiteStorage WebsiteStorage
	ReportStorage  ReportStorage
	Client         *redis.Client
}

func (a *Analyzer) Inspect(w *Website) error {
	ws, err := json.Marshal(w)
	if err != nil {
		return err
	}

	icmd := a.Client.LPush("pending:websites", string(ws))
	if icmd.Err() != nil {
		return icmd.Err()
	}

	return nil
}

func (a *Analyzer) StartReporting(cb func(*Website)) {
	for {
		ss, err := a.Client.BRPop(time.Second*5, "inspect:results").Result()
		if err != nil {
			if err != redis.Nil {
				fmt.Println("Error reading redis list: ", err)
			}
			continue
		}

		if ss != nil {
			var result = &Result{}
			err = json.Unmarshal([]byte(ss[1]), result)
			if err != nil {
				fmt.Println("Error parsing results from redis: ", err.Error())
				continue
			}

			website, err := a.saveReport(result)
			if err != nil {
				fmt.Println("Error saving report: ", err.Error())
				continue
			}
			cb(website)
		}
	}
}

func (a *Analyzer) saveReport(res *Result) (*Website, error) {
	website, err := a.WebsiteStorage.Get(res.WebsiteID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	report := &Report{
		UserID:          website.UserID,
		WebsiteID:       website.ID,
		StartedIn:       res.StartedIn,
		LoadedIn:        res.LoadedIn,
		ResourceCheckIn: res.ResourceCheckIn,
		HTMLCheckIn:     res.HTMLCheckIn,
		TotalIn:         res.TotalIn,
		CreatedAt:       &now,
	}

	err = a.ReportStorage.Save(report)
	if err != nil {
		return nil, err
	}

	for _, match := range res.Matches {
		match.ReportID = report.ID
		match.WebsiteID = website.ID
	}
	website.Matches = res.Matches
	website.InspectedAt = &now

	err = a.WebsiteStorage.Save(website)
	if err != nil {
		return nil, err
	}

	err = a.WebsiteStorage.AddTags([]*Website{website})
	if err != nil {
		return nil, err
	}

	return website, nil
}
