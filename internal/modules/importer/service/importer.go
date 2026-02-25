package service

import (
	"context"
	"time"
)

type ImportItem struct {
	ExternalID   string
	Title        string
	Description  string
	ProgramType  string
	Duration     string
	Thumbnail    string
	VideoURL     string
	PublishedAt  *time.Time
	LanguageCode string
	CategorySlug string
}

type Importer interface {
	SourceType() string
	Fetch(ctx context.Context, baseURL string, since *time.Time) ([]ImportItem, error)
}
