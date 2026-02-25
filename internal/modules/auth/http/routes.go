package http

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r *chi.Mux, h *Handler) {
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Post("/refresh", h.Refresh)
		r.Post("/logout", h.Logout)
	})
}
