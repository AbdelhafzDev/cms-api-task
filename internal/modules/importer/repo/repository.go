package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"cms-api/internal/modules/importer/entity"
	"cms-api/internal/pkg/apperror"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) ListSources(ctx context.Context) ([]*entity.ImportSource, error) {
	var sources []*entity.ImportSource
	if err := r.db.SelectContext(ctx, &sources, queryListSources); err != nil {
		return nil, err
	}
	return sources, nil
}

func (r *repository) GetSourceByID(ctx context.Context, id int64) (*entity.ImportSource, error) {
	var src entity.ImportSource
	if err := r.db.GetContext(ctx, &src, queryGetSourceByID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return &src, nil
}

func (r *repository) CreateLog(ctx context.Context, log *entity.ImportLog) error {
	_, err := r.db.ExecContext(ctx, queryCreateLog,
		log.ID,
		log.SourceID,
		log.TriggeredBy,
		log.Status,
		log.RecordsImported,
		log.ErrorMessage,
		log.StartedAt,
		log.FinishedAt,
	)
	return err
}

func (r *repository) UpdateLog(ctx context.Context, log *entity.ImportLog) error {
	_, err := r.db.ExecContext(ctx, queryUpdateLog,
		log.Status,
		log.RecordsImported,
		log.ErrorMessage,
		log.FinishedAt,
		log.ID,
	)
	return err
}
