package interceptor

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"cms-api/internal/pkg/contextutil"
)

type TokenValidator interface {
	ValidateToken(token string) (*TokenClaims, error)
}

type TokenClaims struct {
	UserID    string
	Email     string
	Roles     []string
	SessionID string
}

func UnaryAuth(validator TokenValidator, excludedMethods []string) grpc.UnaryServerInterceptor {
	excluded := make(map[string]struct{})
	for _, method := range excludedMethods {
		excluded[method] = struct{}{}
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if _, ok := excluded[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		claims, err := extractAndValidateToken(ctx, validator)
		if err != nil {
			return nil, err
		}

		ctx = contextutil.WithUserID(ctx, claims.UserID)
		ctx = contextutil.WithEmail(ctx, claims.Email)
		ctx = contextutil.WithRoles(ctx, claims.Roles)
		ctx = contextutil.WithSessionID(ctx, claims.SessionID)

		return handler(ctx, req)
	}
}

func StreamAuth(validator TokenValidator, excludedMethods []string) grpc.StreamServerInterceptor {
	excluded := make(map[string]struct{})
	for _, method := range excludedMethods {
		excluded[method] = struct{}{}
	}

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if _, ok := excluded[info.FullMethod]; ok {
			return handler(srv, ss)
		}

		claims, err := extractAndValidateToken(ss.Context(), validator)
		if err != nil {
			return err
		}

		wrappedStream := &authenticatedStream{
			ServerStream: ss,
			ctx:          contextutil.WithUserID(ss.Context(), claims.UserID),
		}

		return handler(srv, wrappedStream)
	}
}

func extractAndValidateToken(ctx context.Context, validator TokenValidator) (*TokenClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing authorization header")
	}

	authHeader := authHeaders[0]
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, status.Error(codes.Unauthenticated, "invalid authorization header format")
	}

	token := parts[1]
	claims, err := validator.ValidateToken(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	return claims, nil
}

type authenticatedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *authenticatedStream) Context() context.Context {
	return s.ctx
}
