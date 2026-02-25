package http

import (
	"github.com/go-chi/chi/v5"

	"cms-api/internal/transport/http/middleware"
)

func RegisterRoutes(r *chi.Mux, auth *middleware.AuthMiddleware, h *Handler) {
	r.Route("/api/v1/programs", func(r chi.Router) {
		r.Use(auth.Middleware)

		r.With(middleware.RequireRole("admin", "editor")).Get("/", h.List)
		r.With(middleware.RequireRole("admin", "editor")).Get("/{id}", h.GetByID)
		r.With(middleware.RequireRole("admin", "editor")).Post("/", h.Create)
		r.With(middleware.RequireRole("admin", "editor")).Put("/{id}", h.Update)
		r.With(middleware.RequireRole("admin")).Delete("/{id}", h.Delete)
	})
}
