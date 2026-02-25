package entity

import "time"

type ImportSource struct {
	ID         int64     `db:"id"`
	Name       string    `db:"name"`
	SourceType string    `db:"source_type"`
	BaseURL    string    `db:"base_url"`
	IsActive   bool      `db:"is_active"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type ImportLog struct {
	ID              string     `db:"id"`
	SourceID        int64      `db:"source_id"`
	TriggeredBy     *string    `db:"triggered_by"`
	Status          string     `db:"status"`
	RecordsImported int        `db:"records_imported"`
	ErrorMessage    string     `db:"error_message"`
	StartedAt       *time.Time `db:"started_at"`
	FinishedAt      *time.Time `db:"finished_at"`
	CreatedAt       time.Time  `db:"created_at"`
}
