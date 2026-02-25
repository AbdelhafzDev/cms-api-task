package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"

	"cms-api/internal/infra/cache"
	"cms-api/internal/infra/search"
	"cms-api/internal/modules/discovery/dto"
	"cms-api/internal/modules/discovery/entity"
)

type fakeDiscoveryRepo struct {
	mu       sync.Mutex
	listResp []*entity.Program
	getResp  *entity.Program
	listErr  error
	getErr   error
	listHits int
	getHits  int
}

func (f *fakeDiscoveryRepo) List(ctx context.Context, limit int, cursorPublishedAt *time.Time, cursorID string) ([]*entity.Program, error) {
	_ = ctx
	_ = limit
	_ = cursorPublishedAt
	_ = cursorID
	f.mu.Lock()
	defer f.mu.Unlock()
	f.listHits++
	return f.listResp, f.listErr
}

func (f *fakeDiscoveryRepo) GetByID(ctx context.Context, id string) (*entity.Program, error) {
	_ = ctx
	_ = id
	f.mu.Lock()
	defer f.mu.Unlock()
	f.getHits++
	return f.getResp, f.getErr
}

type fakeCache struct {
	mu   sync.Mutex
	data map[string][]byte
}

func newFakeCache() *fakeCache {
	return &fakeCache{data: make(map[string][]byte)}
}

func (c *fakeCache) Get(ctx context.Context, key string) ([]byte, error) {
	_ = ctx
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.data[key]
	if !ok {
		return nil, cache.ErrCacheMiss
	}
	return val, nil
}

func (c *fakeCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	_ = ctx
	_ = ttl
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
	return nil
}

func (c *fakeCache) Delete(ctx context.Context, keys ...string) error {
	_ = ctx
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, k := range keys {
		delete(c.data, k)
	}
	return nil
}

type fakeSearcher struct{}

func (f *fakeSearcher) Search(ctx context.Context, index string, req search.SearchRequest) (*search.SearchResult, error) {
	_ = ctx
	_ = index
	_ = req
	return &search.SearchResult{Hits: []json.RawMessage{}, Page: req.Page, PerPage: req.PerPage, TotalHits: 0}, nil
}

func makeProgram(id string, publishedAt time.Time) *entity.Program {
	return &entity.Program{
		ID:          id,
		Title:       "Test",
		Description: "Desc",
		ProgramType: "podcast",
		Duration:    sql.NullString{String: "01:00:00", Valid: true},
		PublishedAt: sql.NullTime{Time: publishedAt, Valid: true},
		Thumbnail:   "https://example.com/thumb.jpg",
		VideoURL:    "https://example.com/video.mp4",
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func TestDiscoveryService_List_UsesCache(t *testing.T) {
	repo := &fakeDiscoveryRepo{}
	cacheStore := newFakeCache()
	searcher := &fakeSearcher{}
	log := zap.NewNop()

	svc := New(repo, searcher, cacheStore, log)

	p1 := makeProgram("1", time.Now().Add(-time.Hour))
	repo.listResp = []*entity.Program{p1}

	resp1, err := svc.List(context.Background(), "", 20)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if repo.listHits != 1 {
		t.Fatalf("expected repo list hits 1, got %d", repo.listHits)
	}

	p2 := makeProgram("2", time.Now().Add(-2*time.Hour))
	repo.listResp = []*entity.Program{p2}

	resp2, err := svc.List(context.Background(), "", 20)
	if err != nil {
		t.Fatalf("list cached: %v", err)
	}
	if repo.listHits != 1 {
		t.Fatalf("expected repo list hits to remain 1, got %d", repo.listHits)
	}

	if resp1.Items[0].ID != resp2.Items[0].ID {
		t.Fatalf("expected cached response to match first call")
	}
}

func TestDiscoveryService_GetByID_UsesCache(t *testing.T) {
	repo := &fakeDiscoveryRepo{}
	cacheStore := newFakeCache()
	searcher := &fakeSearcher{}
	log := zap.NewNop()

	svc := New(repo, searcher, cacheStore, log)

	p1 := makeProgram("1", time.Now())
	repo.getResp = p1

	resp1, err := svc.GetByID(context.Background(), "1")
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if repo.getHits != 1 {
		t.Fatalf("expected repo get hits 1, got %d", repo.getHits)
	}

	p2 := makeProgram("2", time.Now())
	repo.getResp = p2

	resp2, err := svc.GetByID(context.Background(), "1")
	if err != nil {
		t.Fatalf("get by id cached: %v", err)
	}
	if repo.getHits != 1 {
		t.Fatalf("expected repo get hits to remain 1, got %d", repo.getHits)
	}

	if resp1.ID != resp2.ID {
		t.Fatalf("expected cached response to match first call")
	}
}

func TestDiscoveryService_Search_PassesThrough(t *testing.T) {
	repo := &fakeDiscoveryRepo{}
	cacheStore := newFakeCache()
	searcher := &fakeSearcher{}
	log := zap.NewNop()

	svc := New(repo, searcher, cacheStore, log)

	_, err := svc.Search(context.Background(), &dto.SearchRequest{
		Query:   "test",
		Page:    1,
		PerPage: 10,
	})
	if err != nil {
		t.Fatalf("search: %v", err)
	}
}
