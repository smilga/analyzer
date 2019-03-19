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
}

type Match struct {
	PatternID PatternID
	Value     string
}

type ResponseTime struct {
	Loaded        string
	ResourceCheck string
	HTMLCheck     string
	Total         string
}

type Result struct {
	Time    ResponseTime
	Matches []*Match
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

	// TODO
	// store website, user, duration, how many requests tested etc for debug
	// USER REQUESTS TABLE for pricing, debug etc
	fmt.Println(result.Time)

	now := time.Now()
	w.MatchedPatterns = result.Matches
	w.SearchedAt = &now

	err = a.WebsiteStorage.Save(w)
	if err != nil {
		return err
	}

	return nil
}
