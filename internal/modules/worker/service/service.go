package service

import (
	"context"
	"time"

	"go.uber.org/zap"

	"cms-api/internal/config"
	"cms-api/internal/infra/search"
	"cms-api/internal/modules/worker/repo"
)

const indexName = "programs"

type service struct {
	repo   repo.Repository
	search search.Indexer
	cfg    config.WorkerConfig
	log    *zap.Logger
}

func New(repo repo.Repository, search search.Indexer, cfg *config.Config, log *zap.Logger) Service {
	return &service{
		repo:   repo,
		search: search,
		cfg:    cfg.Worker,
		log:    log.Named("worker"),
	}
}

func (s *service) Start(ctx context.Context) {
	ticker := time.NewTicker(s.cfg.PollInterval)
	defer ticker.Stop()

	s.log.Info("Worker started",
		zap.Duration("poll_interval", s.cfg.PollInterval),
		zap.Int("batch_size", s.cfg.BatchSize),
		zap.Int("max_attempts", s.cfg.MaxAttempts),
	)

	for {
		select {
		case <-ctx.Done():
			s.log.Info("Worker stopped")
			return
		case <-ticker.C:
			s.processBatch(ctx)
		}
	}
}
