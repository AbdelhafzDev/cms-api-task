package dto

import (
	"cms-api/internal/modules/program/entity"
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
		Status:       p.Status,
		CategoryID:   dbutil.NullInt64ToInt64Ptr(p.CategoryID),
		CategoryName: dbutil.NullStringToPtr(p.CategoryName),
		LanguageID:   dbutil.NullInt64ToInt64Ptr(p.LanguageID),
		LanguageCode: dbutil.NullStringToPtr(p.LanguageCode),
		CreatedBy:    dbutil.NullStringToPtr(p.CreatedBy),
		UpdatedBy:    dbutil.NullStringToPtr(p.UpdatedBy),
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
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
