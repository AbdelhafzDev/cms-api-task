package repo

import (
	"context"
	"time"

	"cms-api/internal/modules/program/entity"
)

type Repository interface {
	Create(ctx context.Context, p *entity.Program) error
	Update(ctx context.Context, p *entity.Program) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*entity.Program, error)
	List(ctx context.Context, limit int, cursorCreatedAt *time.Time, cursorID string) ([]*entity.Program, error)
}
