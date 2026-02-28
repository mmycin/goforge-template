package cache

import (
	"context"
	"fmt"
	"time"
)

var (
	Global Cache
	Memory *RistrettoStore
	Redis  *RedisStore
)

// Cache is the global interface for caching
type Cache interface {
	Get(ctx context.Context, key string, val any) error
	Set(ctx context.Context, key string, val any, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Flush(ctx context.Context) error
}

// Convenience wrappers
func Get(ctx context.Context, key string, val any) error {
	if Global == nil {
		return fmt.Errorf("cache not initialized")
	}
	return Global.Get(ctx, key, val)
}

func Set(ctx context.Context, key string, val any, ttl time.Duration) error {
	if Global == nil {
		return fmt.Errorf("cache not initialized")
	}
	return Global.Set(ctx, key, val, ttl)
}

func Delete(ctx context.Context, key string) error {
	if Global == nil {
		return fmt.Errorf("cache not initialized")
	}
	return Global.Delete(ctx, key)
}

func Flush(ctx context.Context) error {
	if Global == nil {
		return fmt.Errorf("cache not initialized")
	}
	return Global.Flush(ctx)
}
