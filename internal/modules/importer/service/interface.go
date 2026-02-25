package service

import (
	"context"

	"cms-api/internal/modules/importer/dto"
)

type Service interface {
	ListSources(ctx context.Context) (*dto.ImportSourceListResponse, error)
	RunSource(ctx context.Context, sourceID int64) (*dto.ImportRunResponse, error)
}
