package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"cms-api/internal/modules/worker/entity"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) ClaimPendingJobs(ctx context.Context, batchSize int) ([]entity.SearchIndexJob, error) {
	var jobs []entity.SearchIndexJob
	err := r.db.SelectContext(ctx, &jobs, queryClaimPendingJobs, batchSize)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *repository) MarkCompleted(ctx context.Context, jobID string) error {
	_, err := r.db.ExecContext(ctx, queryMarkCompleted, jobID)
	return err
}

func (r *repository) MarkFailed(ctx context.Context, jobID string, errMsg string, nextSchedule time.Time) error {
	_, err := r.db.ExecContext(ctx, queryMarkFailed, jobID, errMsg, nextSchedule)
	return err
}

func (r *repository) MarkDead(ctx context.Context, jobID string, errMsg string) error {
	_, err := r.db.ExecContext(ctx, queryMarkDead, jobID, errMsg)
	return err
}

func (r *repository) GetProgramForIndex(ctx context.Context, programID string) (*entity.ProgramDocument, error) {
	row := r.db.QueryRowContext(ctx, queryGetProgramForIndex, programID)

	var doc entity.ProgramDocument
	var duration, publishedAt, category, language sql.NullString
	var createdAt time.Time

	err := row.Scan(
		&doc.ID,
		&doc.Title,
		&doc.Description,
		&doc.ProgramType,
		&doc.Status,
		&duration,
		&publishedAt,
		&category,
		&language,
		&doc.Thumbnail,
		&doc.VideoURL,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}

	if duration.Valid {
		doc.Duration = &duration.String
	}
	if publishedAt.Valid {
		doc.PublishedAt = &publishedAt.String
	}
	if category.Valid {
		doc.Category = &category.String
	}
	if language.Valid {
		doc.Language = &language.String
	}
	doc.CreatedAt = createdAt.Format(time.RFC3339)

	return &doc, nil
}
