package search

import (
	"context"
	"encoding/json"
)

type SearchResult struct {
	Hits      []json.RawMessage
	Page      int
	PerPage   int
	TotalHits int64
}

type SearchRequest struct {
	Query   string
	Filter  string
	Sort    []string
	Page    int
	PerPage int
}

type IndexConfig struct {
	SearchableAttributes []string
	FilterableAttributes []string
	SortableAttributes   []string
}

type Searcher interface {
	Search(ctx context.Context, index string, req SearchRequest) (*SearchResult, error)
}

type Indexer interface {
	EnsureIndex(ctx context.Context, index string, primaryKey string, cfg IndexConfig) error
	AddDocuments(ctx context.Context, index string, docs []any) error
	DeleteDocument(ctx context.Context, index string, docID string) error
}
