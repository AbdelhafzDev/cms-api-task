package grpc

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"cms-api/internal/config"
)

type Server struct {
	server *grpc.Server
	log    *zap.Logger
	cfg    *config.Config
}

func NewServer(cfg *config.Config, log *zap.Logger) *Server {
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(cfg.GRPC.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(cfg.GRPC.MaxSendMsgSize),
	}

	server := grpc.NewServer(opts...)

	if cfg.IsDevelopment() {
		reflection.Register(server)
	}

	return &Server{
		server: server,
		log:    log,
		cfg:    cfg,
	}
}

func (s *Server) Server() *grpc.Server {
	return s.server
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.GRPC.Host, s.cfg.GRPC.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.log.Info("Starting gRPC server", zap.String("addr", addr))

	if err := s.server.Serve(listener); err != nil {
		return fmt.Errorf("gRPC server error: %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("Shutting down gRPC server")

	stopped := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		s.server.Stop()
		return ctx.Err()
	case <-stopped:
		return nil
	}
}

func RegisterLifecycle(lc fx.Lifecycle, server *Server, cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.Start(); err != nil {
					server.log.Error("gRPC server failed", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			shutdownCtx, cancel := context.WithTimeout(ctx, cfg.GRPC.ShutdownTimeout)
			defer cancel()
			return server.Stop(shutdownCtx)
		},
	})
}
