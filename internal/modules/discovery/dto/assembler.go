package dto

import (
	"encoding/json"

	"cms-api/internal/modules/discovery/entity"
	"cms-api/internal/pkg/dbutil"
)

func ToResponse(p *entity.Program) *ProgramResponse {
	resp := &ProgramResponse{
		ID:           p.ID,
		Title:        p.Title,
		Description:  p.Description,
		ProgramType:  p.ProgramType,
		Duration:     dbutil.NullStringToPtr(p.Duration),
		Thumbnail:    p.Thumbnail,
		VideoURL:     p.VideoURL,
		CategoryName: dbutil.NullStringToPtr(p.CategoryName),
		LanguageCode: dbutil.NullStringToPtr(p.LanguageCode),
	}

	if p.PublishedAt.Valid {
		resp.PublishedAt = &p.PublishedAt.Time
	}

	return resp
}

func ToListResponse(programs []*entity.Program, nextCursor string, hasNext bool) *ProgramListResponse {
	items := make([]*ProgramResponse, 0, len(programs))
	for _, p := range programs {
		items = append(items, ToResponse(p))
	}

	return &ProgramListResponse{
		Items:      items,
		NextCursor: nextCursor,
		HasNext:    hasNext,
	}
}

func HitsToSearchResponse(hits []json.RawMessage, query string, page, perPage int, totalHits int64) (*SearchResultResponse, error) {
	items := make([]*SearchProgramResponse, 0, len(hits))
	for _, raw := range hits {
		var doc programDocument
		if err := json.Unmarshal(raw, &doc); err != nil {
			return nil, err
		}
		items = append(items, &SearchProgramResponse{
			ID:          doc.ID,
			Title:       doc.Title,
			Description: doc.Description,
			ProgramType: doc.ProgramType,
			Duration:    doc.Duration,
			PublishedAt: doc.PublishedAt,
			Category:    doc.Category,
			Language:    doc.Language,
			Thumbnail:   doc.Thumbnail,
			VideoURL:    doc.VideoURL,
		})
	}

	return &SearchResultResponse{
		Items:          items,
		Query:          query,
		Page:           page,
		PerPage:        perPage,
		EstimatedTotal: totalHits,
	}, nil
}

type programDocument struct {
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
}
