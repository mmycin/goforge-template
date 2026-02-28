package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(host string, port int, password string, db int) *RedisStore {
	return &RedisStore{
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", host, port),
			Password: password,
			DB:       db,
		}),
	}
}

func (s *RedisStore) Get(ctx context.Context, key string, val any) error {
	b, err := s.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return fmt.Errorf("key not found: %s", key)
	} else if err != nil {
		return err
	}
	return json.Unmarshal(b, val)
}

func (s *RedisStore) Set(ctx context.Context, key string, val any, ttl time.Duration) error {
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, key, b, ttl).Err()
}

func (s *RedisStore) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}

func (s *RedisStore) Flush(ctx context.Context) error {
	return s.client.FlushDB(ctx).Err()
}
