package service

import (
	"context"

	"cms-api/internal/modules/program/dto"
)

type Service interface {
	Create(ctx context.Context, req *dto.CreateProgramRequest) (*dto.ProgramResponse, error)
	Update(ctx context.Context, id string, req *dto.UpdateProgramRequest) (*dto.ProgramResponse, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*dto.ProgramResponse, error)
	List(ctx context.Context, cursorStr string, limit int) (*dto.ProgramListResponse, error)
}
