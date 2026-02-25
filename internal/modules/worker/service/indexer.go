package service

import (
	"context"

	"go.uber.org/zap"

	"cms-api/internal/infra/search"
)

func (s *service) EnsureIndex(ctx context.Context) error {
	if err := s.search.EnsureIndex(ctx, indexName, "id", search.IndexConfig{
		SearchableAttributes: []string{"title", "description"},
		FilterableAttributes: []string{"status", "program_type", "category", "language"},
		SortableAttributes:   []string{"published_at", "created_at"},
	}); err != nil {
		return err
	}

	s.log.Info("Meilisearch index configured", zap.String("index", indexName))
	return nil
}
