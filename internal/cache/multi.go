package cache

import (
	"context"
	"time"
)

type MultiStore struct {
	l1 Cache
	l2 Cache
}

func NewMultiStore(l1, l2 Cache) *MultiStore {
	return &MultiStore{l1: l1, l2: l2}
}

func (s *MultiStore) Get(ctx context.Context, key string, val any) error {
	if err := s.l1.Get(ctx, key, val); err == nil {
		return nil
	}
	if err := s.l2.Get(ctx, key, val); err != nil {
		return err
	}
	_ = s.l1.Set(ctx, key, val, 10*time.Minute)
	return nil
}

func (s *MultiStore) Set(ctx context.Context, key string, val any, ttl time.Duration) error {
	_ = s.l1.Set(ctx, key, val, ttl)
	return s.l2.Set(ctx, key, val, ttl)
}

func (s *MultiStore) Delete(ctx context.Context, key string) error {
	_ = s.l1.Delete(ctx, key)
	return s.l2.Delete(ctx, key)
}

func (s *MultiStore) Flush(ctx context.Context) error {
	_ = s.l1.Flush(ctx)
	return s.l2.Flush(ctx)
}
