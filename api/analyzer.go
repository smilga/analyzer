package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"

	uuid "github.com/satori/go.uuid"
)

var script = "puppeteer/index.js"

type Analyzer struct {
	ResultStorage  ResultStorage
	ServiceStorage ServiceStorage
}

func (a *Analyzer) Inspect(w *Website, services []*Service) (*ShortReport, error) {
	s, err := json.Marshal(services)
	if err != nil {
		return nil, err
	}

	var stdOut, stdErr bytes.Buffer
	cmd := exec.Command("node", script, w.URL, string(s))
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	var result = &Result{}
	err = json.Unmarshal(stdOut.Bytes(), result)
	if err != nil {
		return nil, err
	}
	result.WebsiteID = w.ID

	err = a.ResultStorage.Save(result)
	if err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, len(result.DetectedServices))
	for i, ds := range result.DetectedServices {
		ids[i] = ds.ServiceID
	}

	userServices, err := a.ServiceStorage.ManyByUser(w.UserID, ids)
	if err != nil {
		return nil, err
	}

	// TODO Add website tags

	foundServices := make([]*ServiceIdentity, len(userServices))
	for i, us := range userServices {
		foundServices[i] = us.ServiceIdentity
	}

	return &ShortReport{
		WebsiteID:  w.ID,
		UserID:     w.UserID,
		SearchedAt: result.CreatedAt,
		Services:   foundServices,
	}, nil
}

func (a *Analyzer) Report(w *Website) (*Report, error) {
	report := &Report{
		WebsiteID:  w.ID,
		WebsiteURL: w.URL,
		UserID:     w.UserID,
	}

	res, err := a.ResultStorage.LatestByWebsite(w.ID)
	if err != nil {
		return nil, err
	}

	report.Duration = res.Duration

	services, err := a.ServiceStorage.ManyByUser(w.UserID, res.ListServiceIDs())
	if err != nil {
		return nil, err
	}

	for _, ds := range res.DetectedServices {
		us, err := extractService(ds.ServiceID, services)
		if err != nil {
			return nil, err
		}

		match := make([]*MatchReport, len(ds.Match))
		for i, m := range ds.Match {
			pattern, err := us.Pattern(m.PatternID)
			if err != nil {
				fmt.Println("Error getting pattern")
				continue
			}
			match[i] = &MatchReport{
				Pattern: pattern,
				Value:   m.Value,
			}
		}

		report.ServiceReports = append(report.ServiceReports, &ServiceReport{
			ServiceIdentity: us.ServiceIdentity,
			Match:           match,
		})
	}

	return report, nil
}

func extractService(id uuid.UUID, services []*Service) (*Service, error) {
	for _, s := range services {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, errors.New("Error extracting service, service not found")
}
