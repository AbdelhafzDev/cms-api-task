package middleware

import (
	"crypto/rsa"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"cms-api/internal/pkg/contextutil"
	"cms-api/internal/pkg/crypto"
	"cms-api/internal/pkg/httputil"
)

type AuthMiddleware struct {
	publicKey *rsa.PublicKey
	logger    *zap.Logger
}

// If publicKeyPath is empty, the middleware rejects all authenticated requests.
func NewAuthMiddleware(publicKeyPath string, logger *zap.Logger) (*AuthMiddleware, error) {
	if publicKeyPath == "" {
		logger.Warn("Auth middleware initialized without public key — authenticated endpoints will reject all requests")
		return &AuthMiddleware{logger: logger}, nil
	}

	publicKey, err := crypto.LoadPublicKey(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load public key: %w", err)
	}

	logger.Info("Auth middleware initialized", zap.String("public_key_path", publicKeyPath))

	return &AuthMiddleware{
		publicKey: publicKey,
		logger:    logger,
	}, nil
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.publicKey == nil {
			httputil.Unauthorized(w, "authentication not configured")
			return
		}

		tokenString, err := crypto.ExtractBearerToken(r.Header.Get("Authorization"))
		if err != nil {
			m.logger.Debug("Failed to extract bearer token", zap.Error(err))
			httputil.Unauthorized(w, err.Error())
			return
		}

		claims, err := crypto.ValidateToken(m.publicKey, tokenString)
		if err != nil {
			m.logger.Warn("Token validation failed", zap.Error(err))
			httputil.Unauthorized(w, "invalid token")
			return
		}

		userID, ok := crypto.ExtractUserID(claims)
		if !ok {
			m.logger.Warn("Missing or invalid user identifier in token claims",
				zap.Any("claims", claims),
			)
			httputil.Unauthorized(w, "invalid token claims")
			return
		}

		ctx := r.Context()
		ctx = contextutil.WithUserID(ctx, userID)
		ctx = contextutil.WithEmail(ctx, crypto.ExtractEmail(claims))
		ctx = contextutil.WithRoles(ctx, crypto.ExtractRoles(claims))
		ctx = contextutil.WithSessionID(ctx, crypto.ExtractSessionID(claims))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) OptionalMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := crypto.ExtractBearerToken(r.Header.Get("Authorization"))
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		claims, err := crypto.ValidateToken(m.publicKey, tokenString)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		userID, ok := crypto.ExtractUserID(claims)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		ctx = contextutil.WithUserID(ctx, userID)
		ctx = contextutil.WithEmail(ctx, crypto.ExtractEmail(claims))
		ctx = contextutil.WithRoles(ctx, crypto.ExtractRoles(claims))
		ctx = contextutil.WithSessionID(ctx, crypto.ExtractSessionID(claims))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRole(roles ...string) func(next http.Handler) http.Handler {
	roleSet := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		roleSet[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRoles := contextutil.GetRoles(r.Context())
			if len(userRoles) == 0 {
				httputil.Forbidden(w, "no role found")
				return
			}

			for _, ur := range userRoles {
				if _, ok := roleSet[ur]; ok {
					next.ServeHTTP(w, r)
					return
				}
			}

			httputil.Forbidden(w, "insufficient permissions")
		})
	}
}
