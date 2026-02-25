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
	CategoryName *string    `json:"category_name"`
	LanguageCode *string    `json:"language_code"`
}

type ProgramListResponse struct {
	Items      []*ProgramResponse `json:"items"`
	NextCursor string             `json:"next_cursor,omitempty"`
	HasNext    bool               `json:"has_next"`
}

type SearchResultResponse struct {
	Items          []*SearchProgramResponse `json:"items"`
	Query          string                   `json:"query"`
	Page           int                      `json:"page"`
	PerPage        int                      `json:"per_page"`
	EstimatedTotal int64                    `json:"estimated_total"`
}

type SearchProgramResponse struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ProgramType string  `json:"program_type"`
	Duration    *string `json:"duration"`
	PublishedAt *string `json:"published_at"`
	Category    *string `json:"category"`
	Language    *string `json:"language"`
	Thumbnail   string  `json:"thumbnail"`
	VideoURL    string  `json:"video_url"`
}
