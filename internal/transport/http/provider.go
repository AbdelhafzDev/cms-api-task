package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"cms-api/internal/config"
	"cms-api/internal/infra/telemetry"
	"cms-api/internal/transport/http/middleware"
)

type RouterResult struct {
	fx.Out

	Router         *chi.Mux
	AuthMiddleware *middleware.AuthMiddleware
	Handler        http.Handler
}

func NewRouterProvider(params RouterParams) (RouterResult, error) {
	router, authMiddleware, err := NewRouter(params.Config, params.Logger, params.Tracer)
	if err != nil {
		return RouterResult{}, err
	}

	return RouterResult{
		Router:         router,
		AuthMiddleware: authMiddleware,
		Handler:        router,
	}, nil
}

var Module = fx.Module("http",
	fx.Provide(
		func(cfg *config.Config, log *zap.Logger, tracer *telemetry.Tracer) RouterParams {
			return RouterParams{
				Config: cfg,
				Logger: log,
				Tracer: tracer,
			}
		},
		NewRouterProvider,
		NewServer,
	),
	fx.Invoke(RegisterLifecycle),
)
