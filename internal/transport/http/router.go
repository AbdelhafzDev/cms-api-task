package http

import (
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	"cms-api/internal/config"
	httpdocs "cms-api/docs/http"
	"cms-api/internal/infra/telemetry"
	"cms-api/internal/transport/http/health"
	"cms-api/internal/transport/http/middleware"
	"cms-api/internal/transport/http/swagger"
)

type RouteRegistrar func(r chi.Router)

type RouterParams struct {
	Config     *config.Config
	Logger     *zap.Logger
	Tracer     *telemetry.Tracer
	Registrars []RouteRegistrar
}

func NewRouter(cfg *config.Config, log *zap.Logger, tracer *telemetry.Tracer) (*chi.Mux, *middleware.AuthMiddleware, error) {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)

	if tracer != nil && cfg.Telemetry.Enabled {
		tracingMiddleware := middleware.NewTracingMiddleware(tracer)
		r.Use(tracingMiddleware.Middleware)
	}

	r.Use(middleware.Logger(log))
	r.Use(middleware.Recoverer(log))
	r.Use(chimiddleware.Compress(5))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.HTTP.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"Link", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	authMiddleware, err := middleware.NewAuthMiddleware(cfg.JWT.PublicKeyPath, log)
	if err != nil {
		return nil, nil, err
	}

	health.RegisterRoutes(r)
	swagger.RegisterRoutes(r, httpdocs.SpecFS, httpdocs.SpecPath)

	return r, authMiddleware, nil
}
