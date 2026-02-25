package http

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"cms-api/internal/modules/program/dto"
	"cms-api/internal/modules/program/service"
	"cms-api/internal/pkg/crypto"
	"cms-api/internal/transport/http/middleware"
)

type fakeProgramService struct {
	listResp *dto.ProgramListResponse
	listErr  error
}

func (f *fakeProgramService) Create(ctx context.Context, req *dto.CreateProgramRequest) (*dto.ProgramResponse, error) {
	return &dto.ProgramResponse{}, nil
}

func (f *fakeProgramService) Update(ctx context.Context, id string, req *dto.UpdateProgramRequest) (*dto.ProgramResponse, error) {
	return &dto.ProgramResponse{}, nil
}

func (f *fakeProgramService) Delete(ctx context.Context, id string) error {
	return nil
}

func (f *fakeProgramService) GetByID(ctx context.Context, id string) (*dto.ProgramResponse, error) {
	return &dto.ProgramResponse{}, nil
}

func (f *fakeProgramService) List(ctx context.Context, cursorStr string, limit int) (*dto.ProgramListResponse, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	if f.listResp != nil {
		return f.listResp, nil
	}
	return &dto.ProgramListResponse{Items: []*dto.ProgramResponse{}, HasNext: false}, nil
}

func generateKeyPair(t *testing.T) (*rsa.PrivateKey, string, func()) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate rsa key: %v", err)
	}

	pubASN1, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("marshal public key: %v", err)
	}

	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubASN1})
	f, err := os.CreateTemp("", "public-*.pem")
	if err != nil {
		t.Fatalf("create temp pub file: %v", err)
	}
	if _, err := f.Write(pubPEM); err != nil {
		f.Close()
		t.Fatalf("write pub key: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close pub key file: %v", err)
	}

	cleanup := func() {
		_ = os.Remove(f.Name())
	}

	return privateKey, f.Name(), cleanup
}

func makeToken(t *testing.T, key *rsa.PrivateKey, roles []string) string {
	t.Helper()

	claims := map[string]interface{}{
		"sub":   "user-1",
		"email": "user@example.com",
		"roles": roles,
	}
	token, err := crypto.GenerateToken(key, claims, 15*time.Minute)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	return token
}

func TestProgramRoutes_AuthRoleChecks(t *testing.T) {
	privateKey, pubPath, cleanup := generateKeyPair(t)
	defer cleanup()

	auth, err := middleware.NewAuthMiddleware(pubPath, zap.NewNop())
	if err != nil {
		t.Fatalf("auth middleware: %v", err)
	}

	h := NewHandler(&fakeProgramService{}, zap.NewNop())
	router := chi.NewRouter()
	RegisterRoutes(router, auth, h)

	req, _ := http.NewRequest(http.MethodGet, "/v1/programs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}

	userToken := makeToken(t, privateKey, []string{"user"})
	req, _ = http.NewRequest(http.MethodGet, "/v1/programs", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}

	adminToken := makeToken(t, privateKey, []string{"admin"})
	req, _ = http.NewRequest(http.MethodGet, "/v1/programs", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

var _ service.Service = (*fakeProgramService)(nil)
