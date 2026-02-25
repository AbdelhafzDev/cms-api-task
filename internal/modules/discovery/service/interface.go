package service

import (
	"context"

	"cms-api/internal/modules/discovery/dto"
)

type Service interface {
	Search(ctx context.Context, req *dto.SearchRequest) (*dto.SearchResultResponse, error)
	List(ctx context.Context, cursorStr string, limit int) (*dto.ProgramListResponse, error)
	GetByID(ctx context.Context, id string) (*dto.ProgramResponse, error)
}
