package repo

import (
	"context"
	"time"

	"cms-api/internal/modules/discovery/entity"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*entity.Program, error)
	List(ctx context.Context, limit int, cursorPublishedAt *time.Time, cursorID string) ([]*entity.Program, error)
}
