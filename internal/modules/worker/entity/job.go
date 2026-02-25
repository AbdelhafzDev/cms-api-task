package entity

import (
	"database/sql"
	"time"
)

type SearchIndexJob struct {
	ID           string       `db:"id"`
	ProgramID    string       `db:"program_id"`
	Action       string       `db:"action"`
	Status       string       `db:"status"`
	Attempts     int          `db:"attempts"`
	MaxAttempts  int          `db:"max_attempts"`
	LastError    sql.NullString `db:"last_error"`
	ScheduledAt  time.Time    `db:"scheduled_at"`
	ProcessedAt  sql.NullTime `db:"processed_at"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at"`
}

type ProgramDocument struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ProgramType string  `json:"program_type"`
	Status      string  `json:"status"`
	Duration    *string `json:"duration,omitempty"`
	PublishedAt *string `json:"published_at,omitempty"`
	Category    *string `json:"category,omitempty"`
	Language    *string `json:"language,omitempty"`
	Thumbnail   string  `json:"thumbnail"`
	VideoURL    string  `json:"video_url"`
	CreatedAt   string  `json:"created_at"`
}
