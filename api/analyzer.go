package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

var script = "puppeteer/index.js"

type Analyzer struct {
	PatternStorage PatternStorage
	WebsiteStorage WebsiteStorage
	ReportStorage  ReportStorage
}

// TODO proccess multiple websites at once!
func (a *Analyzer) Inspect(w *Website) error {
	patterns, err := a.PatternStorage.All()
	if err != nil {
		return err
	}

	patternMap := make(map[PatternID]*Pattern, len(patterns))
	for _, p := range patterns {
		patternMap[p.ID] = p
	}

	s, err := json.Marshal(patterns)
	if err != nil {
		return err
	}

	fmt.Println("======== Pattern string ==========")
	fmt.Println(string(s))
	fmt.Println("======== Pattern string ==========")

	var stdOut, stdErr bytes.Buffer
	cmd := exec.Command("node", script, w.URL, string(s))
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr
	err = cmd.Run()
	if err != nil {
		return err
	}

	var result = &Result{}
	err = json.Unmarshal(stdOut.Bytes(), result)
	if err != nil {
		return err
	}

	report, err := a.saveReport(result, w)
	if err != nil {
		return nil
	}

	for _, match := range result.Matches {
		match.ReportID = report.ID
		match.WebsiteID = w.ID
	}
	now := time.Now()
	w.Matches = result.Matches
	w.InspectedAt = &now

	err = a.WebsiteStorage.Save(w)
	if err != nil {
		return err
	}

	return nil
}

func (a *Analyzer) saveReport(res *Result, w *Website) (*Report, error) {
	now := time.Now()
	report := &Report{
		UserID:          w.UserID,
		WebsiteID:       w.ID,
		LoadedIn:        res.LoadedIn,
		ResourceCheckIn: res.ResourceCheckIn,
		HTMLCheckIn:     res.HTMLCheckIn,
		TotalIn:         res.TotalIn,
		CreatedAt:       &now,
	}

	err := a.ReportStorage.Save(report)
	if err != nil {
		return nil, err
	}
	return report, nil
}
