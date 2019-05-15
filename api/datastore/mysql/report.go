package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/smilga/analyzer/api"
)

type ReportStore struct {
	DB *sql.DB
}

func (s *ReportStore) Save(r *api.Report) error {
	now := time.Now()
	if r.ID == 0 {
		r.CreatedAt = &now
	}

	// NOTE there is trigger that moves deleted reports to reports_archive table
	_, err := s.DB.Exec(`DELETE from reports where website_id=?`, r.WebsiteID)
	if err != nil {
		return err
	}

	res, err := s.DB.Exec(`
		INSERT INTO reports
		(id, user_id, website_id, started_in, loaded_in, resource_check_in, html_check_in, total_in, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, r.ID, r.UserID, r.WebsiteID, r.StartedIn, r.LoadedIn, r.ResourceCheckIn, r.HTMLCheckIn, r.TotalIn, r.CreatedAt)

	if err != nil {
		return err
	}

	if r.ID == 0 {
		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		r.ID = api.ReportID(id)
	}

	return nil
}

func (s *ReportStore) ByWebsite(id api.WebsiteID) (*api.Report, error) {
	r := &api.Report{}

	err := s.DB.QueryRow(`
		SELECT r.*, w.url FROM reports r
		LEFT JOIN websites w on w.id = r.website_id
		WHERE r.website_id = ?
	`, id).Scan(&r.ID, &r.UserID, &r.WebsiteID, &r.StartedIn, &r.LoadedIn, &r.ResourceCheckIn, &r.HTMLCheckIn, &r.TotalIn, &r.CreatedAt,
		&r.WebsiteURL)
	if err != nil {
		return nil, err
	}

	rows, err := s.DB.Query(`
		SELECT * FROM matches where report_id = ?
	`, r.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	patternIDs := []api.PatternID{}
	for rows.Next() {
		m := &api.Match{}
		err := rows.Scan(&m.ID, &m.PatternID, &m.WebsiteID, &m.ReportID, &m.Value, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		r.Matches = append(r.Matches, m)
		patternIDs = append(patternIDs, m.PatternID)
	}

	if len(patternIDs) == 0 {
		return r, nil
	}

	rows, err = s.DB.Query(fmt.Sprintf("SELECT * FROM patterns WHERE id IN (%s);"), patternIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	patterns := []*api.Pattern{}
	for rows.Next() {
		var p api.Pattern
		err := rows.Scan(&p.ID, &p.Type, &p.Value, &p.Description, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt)
		if err != nil {
			return nil, err
		}
		patterns = append(patterns, &p)
	}

	patternMap := make(map[api.PatternID]*api.Pattern, len(patterns))
	for _, p := range patterns {
		patternMap[p.ID] = p
	}

	for _, m := range r.Matches {
		p, ok := patternMap[m.PatternID]
		if !ok {
			return nil, err
		}
		m.Pattern = p
	}

	return r, nil
}

func NewReportStore(DB *sql.DB) *ReportStore {
	return &ReportStore{DB}
}
