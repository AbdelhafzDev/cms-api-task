package dto

import "time"

type ProgramResponse struct {
	ID           string     `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	ProgramType  string     `json:"program_type"`
	Duration     *string    `json:"duration"`
	PublishedAt  *time.Time `json:"published_at"`
	Thumbnail    string     `json:"thumbnail"`
	VideoURL     string     `json:"video_url"`
	Status       string     `json:"status"`
	CategoryID   *int64     `json:"category_id"`
	CategoryName *string    `json:"category_name"`
	LanguageID   *int64     `json:"language_id"`
	LanguageCode *string    `json:"language_code"`
	CreatedBy    *string    `json:"created_by"`
	UpdatedBy    *string    `json:"updated_by"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type ProgramListResponse struct {
	Items      []*ProgramResponse `json:"items"`
	NextCursor string             `json:"next_cursor,omitempty"`
	HasNext    bool               `json:"has_next"`
}
