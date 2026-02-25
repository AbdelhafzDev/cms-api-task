package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"cms-api/internal/config"
)

var Module = fx.Module("cache",
	fx.Provide(NewRedis),
)

type redisCache struct {
	client *redis.Client
}

func NewRedis(lc fx.Lifecycle, cfg *config.Config, log *zap.Logger) (Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Cache.Addr(),
		Password: cfg.Cache.Password,
		DB:       cfg.Cache.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Info("Redis connected",
		zap.String("host", cfg.Cache.Host),
		zap.Int("port", cfg.Cache.Port),
		zap.Int("db", cfg.Cache.DB),
	)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Info("Closing Redis connection")
			return client.Close()
		},
	})

	return &redisCache{client: client}, nil
}

func (c *redisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrCacheMiss
		}
		return nil, err
	}
	return val, nil
}

func (c *redisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *redisCache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return c.client.Del(ctx, keys...).Err()
}
