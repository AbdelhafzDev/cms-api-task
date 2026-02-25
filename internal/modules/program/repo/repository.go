package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"cms-api/internal/modules/program/entity"
	"cms-api/internal/pkg/apperror"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, p *entity.Program) error {
	_, err := r.db.ExecContext(ctx, queryCreate,
		p.ID, p.Title, p.Description, p.ProgramType, p.Duration,
		p.Thumbnail, p.VideoURL, p.Status, p.CategoryID, p.LanguageID,
		p.CreatedBy, p.UpdatedBy,
	)
	return err
}

func (r *repository) Update(ctx context.Context, p *entity.Program) error {
	result, err := r.db.ExecContext(ctx, queryUpdate,
		p.Title, p.Description, p.ProgramType, p.Duration,
		p.Thumbnail, p.VideoURL, p.Status, p.CategoryID, p.LanguageID,
		p.UpdatedBy, p.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return apperror.ErrNotFound
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, queryDelete, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return apperror.ErrNotFound
	}

	return nil
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

func (r *repository) List(ctx context.Context, limit int, cursorCreatedAt *time.Time, cursorID string) ([]*entity.Program, error) {
	var programs []*entity.Program
	var err error

	if cursorCreatedAt != nil {
		err = r.db.SelectContext(ctx, &programs, queryListAfterCursor, limit, *cursorCreatedAt, cursorID)
	} else {
		err = r.db.SelectContext(ctx, &programs, queryListFirst, limit)
	}

	if err != nil {
		return nil, err
	}
	return programs, nil
}
