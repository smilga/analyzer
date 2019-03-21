package api

import "time"

type MatchID int64

type Match struct {
	ID        MatchID    `db:"id"`
	PatternID PatternID  `db:"pattern_id"`
	Value     string     `db:"value"`
	CreatedAt *time.Time `db:"created_at"`
}
