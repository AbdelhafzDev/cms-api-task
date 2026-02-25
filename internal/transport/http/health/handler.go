package health

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)


func RegisterRoutes(r chi.Router) {
	r.Get("/api/v1/health", healthHandler)
	r.Get("/api/v1/ready", readyHandler)
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"healthy"}`))
}

func readyHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ready"}`))
}
