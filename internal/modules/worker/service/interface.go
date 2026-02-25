package service

import "context"

type Service interface {
	EnsureIndex(ctx context.Context) error
	Start(ctx context.Context)
}
