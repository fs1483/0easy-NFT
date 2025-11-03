package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// New creates a new redis client with sane defaults.
func New(addr, password string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 2,
	})
}

// Ping performs a health check against redis.
func Ping(ctx context.Context, client *redis.Client) error {
	return client.Ping(ctx).Err()
}
