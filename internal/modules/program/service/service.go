package service

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"cms-api/internal/modules/program/dto"
	"cms-api/internal/modules/program/entity"
	"cms-api/internal/modules/program/repo"
	"cms-api/internal/pkg/apperror"
	"cms-api/internal/pkg/contextutil"
	"cms-api/internal/pkg/cursor"
	"cms-api/internal/pkg/dbutil"
	"cms-api/internal/pkg/uuidutil"
)

type service struct {
	repo repo.Repository
	log  *zap.Logger
}

func New(repo repo.Repository, log *zap.Logger) Service {
	return &service{repo: repo, log: log}
}

func (s *service) Create(ctx context.Context, req *dto.CreateProgramRequest) (*dto.ProgramResponse, error) {
	id, err := uuidutil.NewV7String()
	if err != nil {
		return nil, fmt.Errorf("generate uuid: %w", err)
	}

	status := req.Status
	if status == "" {
		status = "active"
	}

	userID := contextutil.GetUserID(ctx)

	p := &entity.Program{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		ProgramType: req.ProgramType,
		Duration:    dbutil.NewNullString(req.Duration),
		Thumbnail:   req.Thumbnail,
		VideoURL:    req.VideoURL,
		Status:      status,
		CreatedBy:   dbutil.NewNullString(userID),
		UpdatedBy:   dbutil.NewNullString(userID),
	}

	if req.CategoryID != nil {
		p.CategoryID = dbutil.NewNullInt64(*req.CategoryID, true)
	}
	if req.LanguageID != nil {
		p.LanguageID = dbutil.NewNullInt64(*req.LanguageID, true)
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("create program: %w", err)
	}

	created, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get created program: %w", err)
	}

	return dto.ToResponse(created), nil
}

func (s *service) Update(ctx context.Context, id string, req *dto.UpdateProgramRequest) (*dto.ProgramResponse, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.ProgramType != nil {
		existing.ProgramType = *req.ProgramType
	}
	if req.Duration != nil {
		existing.Duration = dbutil.NewNullString(*req.Duration)
	}
	if req.Thumbnail != nil {
		existing.Thumbnail = *req.Thumbnail
	}
	if req.VideoURL != nil {
		existing.VideoURL = *req.VideoURL
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}
	if req.CategoryID != nil {
		existing.CategoryID = dbutil.NewNullInt64(*req.CategoryID, true)
	}
	if req.LanguageID != nil {
		existing.LanguageID = dbutil.NewNullInt64(*req.LanguageID, true)
	}

	existing.UpdatedBy = dbutil.NewNullString(contextutil.GetUserID(ctx))

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("update program: %w", err)
	}

	updated, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get updated program: %w", err)
	}

	return dto.ToResponse(updated), nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) GetByID(ctx context.Context, id string) (*dto.ProgramResponse, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return dto.ToResponse(p), nil
}

func (s *service) List(ctx context.Context, cursorStr string, limit int) (*dto.ProgramListResponse, error) {
	var cursorTime *time.Time
	var cursorID string

	if cursorStr != "" {
		t, id, err := cursor.DecodePair(cursorStr)
		if err != nil {
			return nil, apperror.ErrBadRequest
		}
		cursorTime = &t
		cursorID = id
	}

	programs, err := s.repo.List(ctx, limit+1, cursorTime, cursorID)
	if err != nil {
		return nil, fmt.Errorf("list programs: %w", err)
	}

	hasNext := len(programs) > limit
	if hasNext {
		programs = programs[:limit]
	}

	var nextCursor string
	if hasNext && len(programs) > 0 {
		last := programs[len(programs)-1]
		nextCursor = cursor.EncodePair(last.CreatedAt, last.ID)
	}

	return dto.ToListResponse(programs, nextCursor, hasNext), nil
}
