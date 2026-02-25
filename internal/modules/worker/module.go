package worker

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"cms-api/internal/modules/worker/repo"
	"cms-api/internal/modules/worker/service"
)

var Module = fx.Module("worker",
	fx.Provide(repo.New),
	fx.Provide(service.New),
	fx.Invoke(startWorker),
)

func startWorker(lc fx.Lifecycle, svc service.Service, log *zap.Logger) {
	var cancel context.CancelFunc

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := svc.EnsureIndex(ctx); err != nil {
				log.Error("Failed to ensure Meilisearch index", zap.Error(err))
				return err
			}

			var workerCtx context.Context
			workerCtx, cancel = context.WithCancel(context.Background())
			go svc.Start(workerCtx)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if cancel != nil {
				cancel()
			}
			return nil
		},
	})
}
