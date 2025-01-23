// internal/cache/redis.go
package cache

import (
	"context"
	"fmt"
	"time"

	"gorbit/internal/config"

	"github.com/go-redis/redis/v8"
)

// RedisClient wraps the redis client with our own methods
type RedisClient struct {
	client *redis.Client
	cfg    *config.Config
}

func (rc *RedisClient) Ping(ctx context.Context) {
	panic("unimplemented")
}

// NewRedisClient creates a new Redis client wrapper
func NewRedisClient(cfg *config.Config) *RedisClient {
	return &RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		}),
		cfg: cfg,
	}
}

// Connect verifies the connection and returns any error
func (rc *RedisClient) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rc.client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("redis connection failed: %w", err)
	}
	return nil
}

// GetClient returns the underlying redis client
func (rc *RedisClient) GetClient() *redis.Client {
	return rc.client
}

// Close gracefully shuts down the client
func (rc *RedisClient) Close() error {
	return rc.client.Close()
}
