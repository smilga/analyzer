package mysql

import (
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx"
	"github.com/smilga/analyzer/api"
)

type ReportStore struct {
	DB *sqlx.DB
}

func (s *ReportStore) Save(r *api.Report) error {
	now := time.Now()
	if r.ID == 0 {
		r.CreatedAt = &now
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

	rows, err := s.DB.Query(`
		SELECT * FROM (SELECT r.*, w.url FROM reports r LEFT JOIN websites w on w.id = r.website_id where w.id = ? order by created_at desc limit 1) r
		LEFT JOIN matches m ON m.report_id = r.id
	`, id)

	if err != nil {
		return nil, err
	}

	patternIDs := []api.PatternID{}
	for rows.Next() {
		m := &api.Match{}
		err := rows.Scan(&r.ID, &r.UserID, &r.WebsiteID, &r.StartedIn, &r.LoadedIn, &r.ResourceCheckIn, &r.HTMLCheckIn, &r.TotalIn, &r.CreatedAt,
			&r.WebsiteURL,
			&m.ID, &m.PatternID, &m.WebsiteID, &m.ReportID, &m.Value, &m.CreatedAt, &m.DeletedAt,
		)

		if err != nil {
			return nil, err
		}
		r.Matches = append(r.Matches, m)
		patternIDs = append(patternIDs, m.PatternID)
	}
	spew.Dump(r.Matches)
	spew.Dump(patternIDs)

	query, args, err := sqlx.In("SELECT * FROM patterns WHERE id IN (?);", patternIDs)
	if err != nil {
		return nil, err
	}

	patterns := []*api.Pattern{}
	err = s.DB.Select(&patterns, query, args...)
	if err != nil {
		return nil, err
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

func NewReportStore(DB *sqlx.DB) *ReportStore {
	return &ReportStore{DB}
}
