package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"cms-api/internal/infra/cache"
	"cms-api/internal/infra/search"
	"cms-api/internal/modules/discovery/dto"
	"cms-api/internal/modules/discovery/repo"
	"cms-api/internal/pkg/apperror"
	"cms-api/internal/pkg/cursor"
)

const indexName = "programs"

const (
	cacheTTLList = 30 * time.Second
	cacheTTLDetail = 60 * time.Second
)

type service struct {
	repo   repo.Repository
	search search.Searcher
	cache  cache.Cache
	log    *zap.Logger
}

func New(repo repo.Repository, search search.Searcher, cache cache.Cache, log *zap.Logger) Service {
	return &service{repo: repo, search: search, cache: cache, log: log}
}

func (s *service) Search(ctx context.Context, req *dto.SearchRequest) (*dto.SearchResultResponse, error) {
	filter := buildFilter(req)

	searchReq := search.SearchRequest{
		Query:   req.Query,
		Page:    req.Page,
		PerPage: req.PerPage,
		Filter:  filter,
		Sort:    []string{"published_at:desc"},
	}

	result, err := s.search.Search(ctx, indexName, searchReq)
	if err != nil {
		return nil, fmt.Errorf("search programs: %w", err)
	}

	resp, err := dto.HitsToSearchResponse(result.Hits, req.Query, req.Page, req.PerPage, result.TotalHits)
	if err != nil {
		return nil, fmt.Errorf("decode search hits: %w", err)
	}

	return resp, nil
}

func (s *service) List(ctx context.Context, cursorStr string, limit int) (*dto.ProgramListResponse, error) {
	cacheKey := fmt.Sprintf("discovery:list:%s:%d", cursorStr, limit)

	if data, err := s.cache.Get(ctx, cacheKey); err == nil {
		var resp dto.ProgramListResponse
		if err := json.Unmarshal(data, &resp); err == nil {
			return &resp, nil
		}
	}

	var cursorTime *time.Time
	var cursorID string

	if cursorStr != "" {
		t, id, err := cursor.DecodePair(cursorStr)
		if err != nil {
			return nil, apperror.ErrBadRequest
		}
		cursorTime = &t
		cursorID = id
	}

	programs, err := s.repo.List(ctx, limit+1, cursorTime, cursorID)
	if err != nil {
		return nil, fmt.Errorf("list programs: %w", err)
	}

	hasNext := len(programs) > limit
	if hasNext {
		programs = programs[:limit]
	}

	var nextCursor string
	if hasNext && len(programs) > 0 {
		last := programs[len(programs)-1]
		nextCursor = cursor.EncodePair(last.PublishedAt.Time, last.ID)
	}

	resp := dto.ToListResponse(programs, nextCursor, hasNext)

	if data, err := json.Marshal(resp); err == nil {
		_ = s.cache.Set(ctx, cacheKey, data, cacheTTLList)
	}

	return resp, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*dto.ProgramResponse, error) {
	cacheKey := "discovery:id:" + id

	if data, err := s.cache.Get(ctx, cacheKey); err == nil {
		var resp dto.ProgramResponse
		if err := json.Unmarshal(data, &resp); err == nil {
			return &resp, nil
		}
	}

	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := dto.ToResponse(p)

	if data, err := json.Marshal(resp); err == nil {
		_ = s.cache.Set(ctx, cacheKey, data, cacheTTLDetail)
	}

	return resp, nil
}

func escapeFilterValue(s string) string {
	return strings.ReplaceAll(s, `'`, `\'`)
}

func buildFilter(req *dto.SearchRequest) string {
	filters := []string{"status = 'active'"}

	if req.ProgramType != "" {
		filters = append(filters, fmt.Sprintf("program_type = '%s'", escapeFilterValue(req.ProgramType)))
	}
	if req.Category != "" {
		filters = append(filters, fmt.Sprintf("category = '%s'", escapeFilterValue(req.Category)))
	}
	if req.Language != "" {
		filters = append(filters, fmt.Sprintf("language = '%s'", escapeFilterValue(req.Language)))
	}

	return strings.Join(filters, " AND ")
}
