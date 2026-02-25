package search

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/meilisearch/meilisearch-go"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"cms-api/internal/config"
)

var Module = fx.Module("search",
	fx.Provide(NewMeilisearch),
)

type meilisearchClient struct {
	client meilisearch.ServiceManager
}

type SearchOut struct {
	fx.Out

	Searcher Searcher
	Indexer  Indexer
}

func NewMeilisearch(lc fx.Lifecycle, cfg *config.Config, log *zap.Logger) (SearchOut, error) {
	client := meilisearch.New(cfg.Search.Addr(), meilisearch.WithAPIKey(cfg.Search.MasterKey))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.HealthWithContext(ctx); err != nil {
		return SearchOut{}, fmt.Errorf("failed to connect to meilisearch: %w", err)
	}

	log.Info("Meilisearch connected",
		zap.String("host", cfg.Search.Host),
		zap.Int("port", cfg.Search.Port),
	)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Info("Meilisearch client stopped")
			return nil
		},
	})

	ms := &meilisearchClient{client: client}
	return SearchOut{
		Searcher: ms,
		Indexer:  ms,
	}, nil
}

func (m *meilisearchClient) Search(ctx context.Context, index string, req SearchRequest) (*SearchResult, error) {
	msReq := &meilisearch.SearchRequest{
		Page:        int64(req.Page),
		HitsPerPage: int64(req.PerPage),
		Filter:      req.Filter,
		Sort:        req.Sort,
	}

	resp, err := m.client.Index(index).SearchWithContext(ctx, req.Query, msReq)
	if err != nil {
		return nil, fmt.Errorf("meilisearch search: %w", err)
	}

	hits := make([]json.RawMessage, 0, len(resp.Hits))
	for _, hit := range resp.Hits {
		raw, err := json.Marshal(hit)
		if err != nil {
			return nil, fmt.Errorf("marshal search hit: %w", err)
		}
		hits = append(hits, raw)
	}

	return &SearchResult{
		Hits:      hits,
		Page:      req.Page,
		PerPage:   req.PerPage,
		TotalHits: resp.TotalHits,
	}, nil
}

func (m *meilisearchClient) EnsureIndex(ctx context.Context, index string, primaryKey string, cfg IndexConfig) error {
	_, err := m.client.GetIndex(index)
	if err != nil {
		_, createErr := m.client.CreateIndex(&meilisearch.IndexConfig{
			Uid:        index,
			PrimaryKey: primaryKey,
		})
		if createErr != nil {
			return fmt.Errorf("create index: %w", createErr)
		}
	}

	idx := m.client.Index(index)

	if _, err := idx.UpdateSearchableAttributes(&cfg.SearchableAttributes); err != nil {
		return fmt.Errorf("update searchable attributes: %w", err)
	}

	filterableAttrs := make([]interface{}, len(cfg.FilterableAttributes))
	for i, attr := range cfg.FilterableAttributes {
		filterableAttrs[i] = attr
	}
	if _, err := idx.UpdateFilterableAttributes(&filterableAttrs); err != nil {
		return fmt.Errorf("update filterable attributes: %w", err)
	}

	if _, err := idx.UpdateSortableAttributes(&cfg.SortableAttributes); err != nil {
		return fmt.Errorf("update sortable attributes: %w", err)
	}

	return nil
}

func (m *meilisearchClient) AddDocuments(ctx context.Context, index string, docs []any) error {
	if _, err := m.client.Index(index).AddDocuments(docs, nil); err != nil {
		return fmt.Errorf("add documents: %w", err)
	}
	return nil
}

func (m *meilisearchClient) DeleteDocument(ctx context.Context, index string, docID string) error {
	if _, err := m.client.Index(index).DeleteDocument(docID, nil); err != nil {
		return fmt.Errorf("delete document: %w", err)
	}
	return nil
}
