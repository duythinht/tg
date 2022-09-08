package model

import (
	"github.com/jackc/pgtype"
)

// Scan model for pg schema
type Scan struct {
	ID           int64            `json:"id"`
	RepositoryID int64            `json:"repository_id"`
	Status       string           `json:"status"`
	Findings     pgtype.JSONB     `json:"findings"`
	QueuedAt     pgtype.Timestamp `json:"queued_at"`
	ScanningAt   pgtype.Timestamp `json:"scanning_at"`
	FinishedAt   pgtype.Timestamp `json:"finished_at"`
}
