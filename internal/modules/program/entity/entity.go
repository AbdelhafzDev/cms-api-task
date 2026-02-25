package entity

import (
	"database/sql"
	"time"
)

type Program struct {
	ID           string         `db:"id"`
	Title        string         `db:"title"`
	Description  string         `db:"description"`
	ProgramType  string         `db:"program_type"`
	Duration     sql.NullString `db:"duration"`
	PublishedAt  sql.NullTime   `db:"published_at"`
	Thumbnail    string         `db:"thumbnail"`
	VideoURL     string         `db:"video_url"`
	ExternalID   sql.NullString `db:"external_id"`
	Status       string         `db:"status"`
	CategoryID   sql.NullInt64  `db:"category_id"`
	LanguageID   sql.NullInt64  `db:"language_id"`
	ImportSource sql.NullInt64  `db:"import_source_id"`
	CreatedBy    sql.NullString `db:"created_by"`
	UpdatedBy    sql.NullString `db:"updated_by"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
	DeletedAt    sql.NullTime   `db:"deleted_at"`

	// Joined fields
	CategoryName sql.NullString `db:"category_name"`
	LanguageCode sql.NullString `db:"language_code"`
}
