package repo

import (
	"context"
	"time"

	"cms-api/internal/modules/worker/entity"
)

type Repository interface {
	ClaimPendingJobs(ctx context.Context, batchSize int) ([]entity.SearchIndexJob, error)
	MarkCompleted(ctx context.Context, jobID string) error
	MarkFailed(ctx context.Context, jobID string, errMsg string, nextSchedule time.Time) error
	MarkDead(ctx context.Context, jobID string, errMsg string) error
	GetProgramForIndex(ctx context.Context, programID string) (*entity.ProgramDocument, error)
}
