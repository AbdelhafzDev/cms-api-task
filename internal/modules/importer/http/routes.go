package http

import (
	"github.com/go-chi/chi/v5"

	"cms-api/internal/transport/http/middleware"
)

func RegisterRoutes(r *chi.Mux, auth *middleware.AuthMiddleware, h *Handler) {
	r.Route("/v1/import-sources", func(r chi.Router) {
		r.Use(auth.Middleware)

		r.With(middleware.RequireRole("admin", "editor")).Get("/", h.ListSources)
		r.With(middleware.RequireRole("admin", "editor")).Post("/{id}/run", h.RunSource)
	})
}
