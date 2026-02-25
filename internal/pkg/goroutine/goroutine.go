package goroutine

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// SafeWithTimeout launches a goroutine with timeout and panic recovery.
// Use for: Quick background tasks (analytics, cleanup, notifications)
func SafeWithTimeout(logger *zap.Logger, d time.Duration, fn func(ctx context.Context)) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("goroutine panic", zap.Any("panic", r), zap.Stack("stacktrace"))
			}
		}()

		ctx, cancel := context.WithTimeout(context.Background(), d)
		defer cancel()

		fn(ctx)
	}()
}

// Background launches a panic-safe goroutine that runs until completion.
// No timeout, no cancellation - the function will run to completion.
// Use for: Critical background tasks that must complete (database writes, webhooks, cleanup)
func Background(logger *zap.Logger, fn func(ctx context.Context)) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("background goroutine panic", zap.Any("panic", r), zap.Stack("stacktrace"))
			}
		}()

		fn(context.Background())
	}()
}

// BackgroundWithCancel launches a panic-safe goroutine with cancellation support.
// The function will be canceled when the cancel channel is closed.
// Use for: Background tasks that should stop gracefully (client disconnect, shutdown, resource cleanup)
func BackgroundWithCancel(logger *zap.Logger, cancel <-chan struct{}, fn func(ctx context.Context)) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("background goroutine panic", zap.Any("panic", r), zap.Stack("stacktrace"))
			}
		}()

		ctx, ctxCancel := context.WithCancel(context.Background())
		defer ctxCancel()

		go func() {
			<-cancel
			ctxCancel()
		}()

		fn(ctx)
	}()
}

// Stream launches a panic-safe goroutine for streaming operations.
// Closes errChan when done, sends any error before closing.
// Use for: Streaming responses (audio, video, SSE)
func Stream(ctx context.Context, logger *zap.Logger, errChan chan<- error, fn func(ctx context.Context) error) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("stream goroutine panic", zap.Any("panic", r), zap.Stack("stacktrace"))
				errChan <- fmt.Errorf("internal error")
			}
			close(errChan)
		}()

		if err := fn(ctx); err != nil {
			errChan <- err
		}
	}()
}
