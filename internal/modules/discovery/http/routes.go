package http

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
)

func RegisterRoutes(r *chi.Mux, h *Handler) {
	r.Route("/api/v1/discover/programs", func(r chi.Router) {
		r.Use(httprate.LimitByIP(100, 1*time.Minute))
		r.Get("/search", h.Search)
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
	})
}
