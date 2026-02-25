package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"cms-api/internal/modules/auth/dto"
	"cms-api/internal/modules/auth/entity"
	"cms-api/internal/modules/auth/repo"
	"cms-api/internal/pkg/apperror"
	"cms-api/internal/pkg/crypto"
)

type service struct {
	repo           repo.Repository
	privateKey     *rsa.PrivateKey
	expiry         time.Duration
	refreshExpiry  time.Duration
	log            *zap.Logger
}

func New(repo repo.Repository, privateKey *rsa.PrivateKey, expiry time.Duration, refreshExpiry time.Duration, log *zap.Logger) Service {
	return &service{
		repo:          repo,
		privateKey:    privateKey,
		expiry:        expiry,
		refreshExpiry: refreshExpiry,
		log:           log,
	}
}

func (s *service) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if !crypto.CheckPassword(req.Password, user.PasswordHash) {
		return nil, apperror.ErrInvalidCredentials
	}

	if !user.IsActive() {
		return nil, apperror.ErrUserInactive
	}

	roles, err := s.repo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return s.issueTokenPair(ctx, user, roles)
}

func (s *service) Refresh(ctx context.Context, req *dto.RefreshRequest) (*dto.LoginResponse, error) {
	tokenHash, err := crypto.HashToken(req.RefreshToken)
	if err != nil {
		return nil, apperror.ErrInvalidToken
	}

	rt, err := s.repo.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	// Reuse detection: if the token was already revoked, revoke ALL user tokens
	if rt.IsRevoked() {
		s.log.Warn("refresh token reuse detected", zap.String("user_id", rt.UserID))
		_ = s.repo.RevokeAllUserRefreshTokens(ctx, rt.UserID)
		return nil, apperror.ErrTokenRevoked
	}

	if rt.IsExpired() {
		return nil, apperror.ErrTokenExpired
	}

	// Revoke the old refresh token (rotation)
	if err := s.repo.RevokeRefreshToken(ctx, rt.ID); err != nil {
		return nil, err
	}

	// Verify user is still active
	user, err := s.repo.GetUserByID(ctx, rt.UserID)
	if err != nil {
		return nil, err
	}

	if !user.IsActive() {
		_ = s.repo.RevokeAllUserRefreshTokens(ctx, user.ID)
		return nil, apperror.ErrUserInactive
	}

	roles, err := s.repo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return s.issueTokenPair(ctx, user, roles)
}

func (s *service) Logout(ctx context.Context, req *dto.LogoutRequest) error {
	tokenHash, err := crypto.HashToken(req.RefreshToken)
	if err != nil {
		return apperror.ErrInvalidToken
	}

	rt, err := s.repo.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return err
	}

	return s.repo.RevokeRefreshToken(ctx, rt.ID)
}

func (s *service) issueTokenPair(ctx context.Context, user *entity.User, roles []string) (*dto.LoginResponse, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"roles": roles,
	}

	accessToken, err := crypto.GenerateToken(s.privateKey, claims, s.expiry)
	if err != nil {
		return nil, err
	}

	rawRefreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	tokenHash, err := crypto.HashToken(rawRefreshToken)
	if err != nil {
		return nil, err
	}

	rt := &entity.RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(s.refreshExpiry),
	}

	if err := s.repo.CreateRefreshToken(ctx, rt); err != nil {
		return nil, err
	}

	return dto.ToLoginResponse(accessToken, s.expiry, rawRefreshToken, s.refreshExpiry), nil
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b), nil
}
