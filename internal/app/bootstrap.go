package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

type GracefulShutdown struct {
	log      *zap.Logger
	shutdown chan struct{}
	done     chan struct{}
}

func NewGracefulShutdown(log *zap.Logger) *GracefulShutdown {
	return &GracefulShutdown{
		log:      log,
		shutdown: make(chan struct{}),
		done:     make(chan struct{}),
	}
}

func (g *GracefulShutdown) Wait() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		g.log.Info("Received shutdown signal", zap.String("signal", sig.String()))
	case <-g.shutdown:
		g.log.Info("Shutdown initiated programmatically")
	}

	close(g.done)
}


func (g *GracefulShutdown) Shutdown() {
	close(g.shutdown)
}

// Done returns a channel that's closed when shutdown is complete
func (g *GracefulShutdown) Done() <-chan struct{} {
	return g.done
}


func RunWithContext(ctx context.Context, log *zap.Logger, fn func() error) error {
	errChan := make(chan error, 1)

	go func() {
		errChan <- fn()
	}()

	select {
	case <-ctx.Done():
		log.Info("Context cancelled, initiating shutdown")
		return ctx.Err()
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("application error: %w", err)
		}
		return nil
	}
}
