package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"cms-api/internal/modules/discovery/entity"
	"cms-api/internal/pkg/apperror"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetByID(ctx context.Context, id string) (*entity.Program, error) {
	var p entity.Program
	if err := r.db.GetContext(ctx, &p, queryGetByID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *repository) List(ctx context.Context, limit int, cursorPublishedAt *time.Time, cursorID string) ([]*entity.Program, error) {
	var programs []*entity.Program
	var err error

	if cursorPublishedAt != nil {
		err = r.db.SelectContext(ctx, &programs, queryListAfterCursor, limit, *cursorPublishedAt, cursorID)
	} else {
		err = r.db.SelectContext(ctx, &programs, queryListFirst, limit)
	}

	if err != nil {
		return nil, err
	}
	return programs, nil
}
