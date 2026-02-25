package repo

import (
	"context"

	"cms-api/internal/modules/importer/entity"
)

type Repository interface {
	ListSources(ctx context.Context) ([]*entity.ImportSource, error)
	GetSourceByID(ctx context.Context, id int64) (*entity.ImportSource, error)

	CreateLog(ctx context.Context, log *entity.ImportLog) error
	UpdateLog(ctx context.Context, log *entity.ImportLog) error
}
