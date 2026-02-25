package telemetry

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"cms-api/internal/config"
)

var Module = fx.Module("telemetry",
	fx.Provide(NewTracerFromConfig),
	fx.Invoke(RegisterLifecycle),
)

func NewTracerFromConfig(cfg *config.Config, log *zap.Logger) (*Tracer, error) {
	telemetryCfg := Config{
		ServiceName:    cfg.App.Name,
		ServiceVersion: cfg.App.Version,
		Environment:    cfg.App.Env,
		OTLPEndpoint:   cfg.Telemetry.OTLPEndpoint,
		Enabled:        cfg.Telemetry.Enabled,
		SampleRate:     cfg.Telemetry.SampleRate,
	}

	tracer, err := NewTracer(context.Background(), telemetryCfg)
	if err != nil {
		return nil, err
	}

	if telemetryCfg.Enabled {
		log.Info("OpenTelemetry tracing enabled",
			zap.String("service", telemetryCfg.ServiceName),
			zap.String("endpoint", telemetryCfg.OTLPEndpoint),
			zap.Float64("sample_rate", telemetryCfg.SampleRate),
		)
	} else {
		log.Info("OpenTelemetry tracing disabled")
	}

	return tracer, nil
}

func RegisterLifecycle(lc fx.Lifecycle, tracer *Tracer, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Info("Shutting down tracer")
			return tracer.Shutdown(ctx)
		},
	})
}
