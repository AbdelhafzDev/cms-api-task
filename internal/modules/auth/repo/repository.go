package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"cms-api/internal/modules/auth/entity"
	"cms-api/internal/pkg/apperror"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.GetContext(ctx, &user, queryGetUserByEmail, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrInvalidCredentials
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	if err := r.db.GetContext(ctx, &user, queryGetUserByID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	var roles []string
	if err := r.db.SelectContext(ctx, &roles, queryGetUserRoles, userID); err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *repository) CreateRefreshToken(ctx context.Context, rt *entity.RefreshToken) error {
	return r.db.QueryRowContext(ctx, queryCreateRefreshToken, rt.UserID, rt.TokenHash, rt.ExpiresAt).
		Scan(&rt.ID, &rt.CreatedAt)
}

func (r *repository) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
	var rt entity.RefreshToken
	if err := r.db.GetContext(ctx, &rt, queryGetRefreshTokenByHash, tokenHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrInvalidToken
		}
		return nil, err
	}
	return &rt, nil
}

func (r *repository) RevokeRefreshToken(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, queryRevokeRefreshToken, id)
	return err
}

func (r *repository) RevokeAllUserRefreshTokens(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, queryRevokeAllUserRefreshTokens, userID)
	return err
}
