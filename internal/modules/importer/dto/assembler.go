package dto

import "cms-api/internal/modules/importer/entity"

func ToSourceResponse(s *entity.ImportSource) *ImportSourceResponse {
	return &ImportSourceResponse{
		ID:         s.ID,
		Name:       s.Name,
		SourceType: s.SourceType,
		BaseURL:    s.BaseURL,
		IsActive:   s.IsActive,
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
	}
}

func ToSourceListResponse(items []*entity.ImportSource) *ImportSourceListResponse {
	resp := &ImportSourceListResponse{
		Items: make([]*ImportSourceResponse, 0, len(items)),
	}
	for _, s := range items {
		resp.Items = append(resp.Items, ToSourceResponse(s))
	}
	return resp
}

func ToRunResponse(log *entity.ImportLog) *ImportRunResponse {
	return &ImportRunResponse{
		LogID:           log.ID,
		SourceID:        log.SourceID,
		Status:          log.Status,
		RecordsImported: log.RecordsImported,
		ErrorMessage:    log.ErrorMessage,
		StartedAt:       log.StartedAt,
		FinishedAt:      log.FinishedAt,
	}
}
