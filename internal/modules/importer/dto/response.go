package dto

import "time"

type ImportSourceResponse struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	SourceType string    `json:"source_type"`
	BaseURL    string    `json:"base_url"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ImportSourceListResponse struct {
	Items []*ImportSourceResponse `json:"items"`
}

type ImportRunResponse struct {
	LogID           string     `json:"log_id"`
	SourceID        int64      `json:"source_id"`
	Status          string     `json:"status"`
	RecordsImported int        `json:"records_imported"`
	ErrorMessage    string     `json:"error_message,omitempty"`
	StartedAt       *time.Time `json:"started_at,omitempty"`
	FinishedAt      *time.Time `json:"finished_at,omitempty"`
}
