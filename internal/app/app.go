package app

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"cms-api/internal/config"
	"cms-api/internal/infra"
	"cms-api/internal/pkg/logger"
	"cms-api/internal/transport"
)


func Run(version string) {
	app := fx.New(

		fx.Supply(version),

		// Core modules
		config.Module,
		logger.Module,
		infra.Module,

		// Feature modules
		FeatureModules,
		transport.Module,


		fx.Invoke(bootstrap),
	)

	app.Run()
}

// bootstrap is called after all dependencies are initialized
func bootstrap(lc fx.Lifecycle, log *zap.Logger, cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Application starting",
				zap.String("env", cfg.App.Env),
				zap.String("name", cfg.App.Name),
			)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Application shutting down")
			return nil
		},
	})
}
