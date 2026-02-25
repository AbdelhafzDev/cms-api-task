package integration

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"

	"cms-api/internal/modules/program/entity"
	"cms-api/internal/modules/program/repo"
	"cms-api/internal/pkg/uuidutil"
)

func openTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		t.Skip("TEST_DB_DSN not set; skipping integration test")
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		t.Fatalf("connect db: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var regclass sql.NullString
	if err := db.GetContext(ctx, &regclass, "SELECT to_regclass('public.programs')"); err != nil {
		db.Close()
		t.Fatalf("check programs table: %v", err)
	}
	if !regclass.Valid {
		db.Close()
		t.Skip("programs table not found; run migrations before tests")
	}

	return db
}

func TestProgramRepository_CRUDAndList(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	repository := repo.New(db)

	id, err := uuidutil.NewV7String()
	if err != nil {
		t.Fatalf("uuid: %v", err)
	}

	p := &entity.Program{
		ID:          id,
		Title:       "Test Program",
		Description: "Test Description",
		ProgramType: "podcast",
		Duration:    sql.NullString{String: "01:00:00", Valid: true},
		Thumbnail:   "https://example.com/thumb.jpg",
		VideoURL:    "https://example.com/video.mp4",
		Status:      "active",
	}

	ctx := context.Background()
	if err := repository.Create(ctx, p); err != nil {
		t.Fatalf("create program: %v", err)
	}
	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(), "DELETE FROM programs WHERE id = $1", id)
	})

	got, err := repository.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if got.Title != p.Title {
		t.Fatalf("unexpected program: got title=%q", got.Title)
	}

	items, err := repository.List(ctx, 10, nil, "")
	if err != nil {
		t.Fatalf("list programs: %v", err)
	}
	if len(items) == 0 {
		t.Fatalf("expected at least one program in list")
	}
}
