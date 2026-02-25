package service

import (
	"context"

	"cms-api/internal/modules/auth/dto"
)

type Service interface {
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	Refresh(ctx context.Context, req *dto.RefreshRequest) (*dto.LoginResponse, error)
	Logout(ctx context.Context, req *dto.LogoutRequest) error
}
