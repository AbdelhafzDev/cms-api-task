package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"cms-api/internal/modules/importer/dto"
	"cms-api/internal/modules/importer/entity"
	"cms-api/internal/modules/importer/repo"
	"cms-api/internal/pkg/apperror"
	"cms-api/internal/pkg/contextutil"
	"cms-api/internal/pkg/uuidutil"
)

type service struct {
	repo     repo.Repository
	registry *Registry
	log      *zap.Logger
}

func New(repo repo.Repository, registry *Registry, log *zap.Logger) Service {
	return &service{
		repo:     repo,
		registry: registry,
		log:      log.Named("importer"),
	}
}

func (s *service) ListSources(ctx context.Context) (*dto.ImportSourceListResponse, error) {
	items, err := s.repo.ListSources(ctx)
	if err != nil {
		return nil, err
	}
	return dto.ToSourceListResponse(items), nil
}

func (s *service) RunSource(ctx context.Context, sourceID int64) (*dto.ImportRunResponse, error) {
	source, err := s.repo.GetSourceByID(ctx, sourceID)
	if err != nil {
		return nil, err
	}
	if !source.IsActive {
		return nil, apperror.NewAppError(apperror.ErrBadRequest, "import source is inactive", http.StatusBadRequest)
	}

	importer := s.registry.Get(source.SourceType)
	if importer == nil {
		return nil, apperror.NewAppError(apperror.ErrServiceUnavailable, "importer not configured for this source type", http.StatusServiceUnavailable)
	}

	logID, err := uuidutil.NewV7String()
	if err != nil {
		return nil, fmt.Errorf("generate log id: %w", err)
	}

	now := time.Now()
	triggeredBy := contextutil.GetUserID(ctx)
	var triggeredByPtr *string
	if triggeredBy != "" {
		triggeredByPtr = &triggeredBy
	}
	log := &entity.ImportLog{
		ID:              logID,
		SourceID:        source.ID,
		TriggeredBy:     triggeredByPtr,
		Status:          "running",
		RecordsImported: 0,
		ErrorMessage:    "",
		StartedAt:       &now,
		FinishedAt:      nil,
	}

	if err := s.repo.CreateLog(ctx, log); err != nil {
		return nil, err
	}

	items, err := importer.Fetch(ctx, source.BaseURL, nil)
	if err != nil {
		log.Status = "failed"
		log.ErrorMessage = err.Error()
		finished := time.Now()
		log.FinishedAt = &finished
		_ = s.repo.UpdateLog(ctx, log)
		return nil, err
	}

	log.Status = "completed"
	log.RecordsImported = len(items)
	finished := time.Now()
	log.FinishedAt = &finished
	if err := s.repo.UpdateLog(ctx, log); err != nil {
		return nil, err
	}

	return dto.ToRunResponse(log), nil
}
