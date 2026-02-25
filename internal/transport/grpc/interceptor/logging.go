package interceptor

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnaryLogging(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		st, _ := status.FromError(err)

		log.Info("gRPC unary request",
			zap.String("method", info.FullMethod),
			zap.String("code", st.Code().String()),
			zap.Duration("duration", duration),
			zap.Error(err),
		)

		return resp, err
	}
}

func StreamLogging(log *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		err := handler(srv, ss)

		duration := time.Since(start)
		st, _ := status.FromError(err)

		log.Info("gRPC stream request",
			zap.String("method", info.FullMethod),
			zap.Bool("client_stream", info.IsClientStream),
			zap.Bool("server_stream", info.IsServerStream),
			zap.String("code", st.Code().String()),
			zap.Duration("duration", duration),
			zap.Error(err),
		)

		return err
	}
}
