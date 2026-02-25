package entity

import (
	"database/sql"
	"time"
)

type Program struct {
	ID          string         `db:"id"`
	Title       string         `db:"title"`
	Description string         `db:"description"`
	ProgramType string         `db:"program_type"`
	Duration    sql.NullString `db:"duration"`
	PublishedAt sql.NullTime   `db:"published_at"`
	Thumbnail   string         `db:"thumbnail"`
	VideoURL    string         `db:"video_url"`
	Status      string         `db:"status"`
	CategoryID  sql.NullInt64  `db:"category_id"`
	LanguageID  sql.NullInt64  `db:"language_id"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`

	// Joined fields
	CategoryName sql.NullString `db:"category_name"`
	LanguageCode sql.NullString `db:"language_code"`
}
