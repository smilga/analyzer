package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
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

type Result struct {
	Duration int64
	Matches  []*Match
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
	fmt.Println(result.Duration)

	matchedPatterns := make([]*MatchedPattern, len(result.Matches))
	for i, m := range result.Matches {
		p, ok := patternMap[m.PatternID]
		if !ok {
			fmt.Println("Error analyzing. Pattern mismatch")
			continue
		}
		matchedPatterns[i] = &MatchedPattern{p, m.Value}
	}

	w.MatchedPatterns = matchedPatterns

	err = a.WebsiteStorage.Save(w)
	if err != nil {
		return err
	}

	return nil
}
