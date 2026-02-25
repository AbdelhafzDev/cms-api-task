package service

import (
	"context"
	"math"
	"time"

	"go.uber.org/zap"

	"cms-api/internal/modules/worker/entity"
)

func (s *service) processBatch(ctx context.Context) {
	jobs, err := s.repo.ClaimPendingJobs(ctx, s.cfg.BatchSize)
	if err != nil {
		s.log.Error("Failed to fetch pending jobs", zap.Error(err))
		return
	}

	for _, job := range jobs {
		if ctx.Err() != nil {
			return
		}
		s.processJob(ctx, job)
	}
}

func (s *service) processJob(ctx context.Context, job entity.SearchIndexJob) {
	if job.Action == "delete" {
		if err := s.search.DeleteDocument(ctx, indexName, job.ProgramID); err != nil {
			s.handleFailure(ctx, job.ID, job.Attempts, err)
			return
		}
		s.markCompleted(ctx, job.ID)
		s.log.Info("Deleted program from index", zap.String("program_id", job.ProgramID))
		return
	}

	doc, err := s.repo.GetProgramForIndex(ctx, job.ProgramID)
	if err != nil {
		s.handleFailure(ctx, job.ID, job.Attempts, err)
		return
	}

	if err := s.search.AddDocuments(ctx, indexName, []any{doc}); err != nil {
		s.handleFailure(ctx, job.ID, job.Attempts, err)
		return
	}

	s.markCompleted(ctx, job.ID)
	s.log.Info("Indexed program", zap.String("program_id", job.ProgramID), zap.String("title", doc.Title))
}

func (s *service) markCompleted(ctx context.Context, jobID string) {
	if err := s.repo.MarkCompleted(ctx, jobID); err != nil {
		s.log.Error("Failed to mark job completed", zap.String("job_id", jobID), zap.Error(err))
	}
}

func (s *service) handleFailure(ctx context.Context, jobID string, currentAttempts int, err error) {
	nextAttempt := currentAttempts + 1

	if nextAttempt >= s.cfg.MaxAttempts {
		if markErr := s.repo.MarkDead(ctx, jobID, err.Error()); markErr != nil {
			s.log.Error("Failed to mark job dead", zap.String("job_id", jobID), zap.Error(markErr))
		}
		s.log.Warn("Job moved to dead letter", zap.String("job_id", jobID), zap.Int("attempts", nextAttempt), zap.Error(err))
		return
	}

	backoff := time.Duration(math.Pow(2, float64(nextAttempt))) * 10 * time.Second
	nextSchedule := time.Now().Add(backoff)

	if markErr := s.repo.MarkFailed(ctx, jobID, err.Error(), nextSchedule); markErr != nil {
		s.log.Error("Failed to mark job failed", zap.String("job_id", jobID), zap.Error(markErr))
	}
	s.log.Warn("Job failed, scheduling retry",
		zap.String("job_id", jobID),
		zap.Int("attempt", nextAttempt),
		zap.Duration("backoff", backoff),
		zap.Error(err),
	)
}
