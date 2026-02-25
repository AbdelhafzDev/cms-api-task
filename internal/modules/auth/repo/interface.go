package repo

import (
	"context"

	"cms-api/internal/modules/auth/entity"
)

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	GetUserRoles(ctx context.Context, userID string) ([]string, error)
	CreateRefreshToken(ctx context.Context, rt *entity.RefreshToken) error
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, id string) error
	RevokeAllUserRefreshTokens(ctx context.Context, userID string) error
}
