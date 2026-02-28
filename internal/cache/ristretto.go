package cache

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/dgraph-io/ristretto"
)

type RistrettoStore struct {
	cache *ristretto.Cache
}

func NewRistrettoStore(maxItems, maxCost int64) (*RistrettoStore, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: maxItems * 10,
		MaxCost:     maxCost,
		BufferItems: 64,
	})
	if err != nil {
		return nil, err
	}
	return &RistrettoStore{cache: cache}, nil
}

func (s *RistrettoStore) Get(ctx context.Context, key string, val any) error {
	v, found := s.cache.Get(key)
	if !found {
		return fmt.Errorf("key not found: %s", key)
	}
	if v == nil {
		return nil
	}

	rv := reflect.ValueOf(val)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("val must be a non-nil pointer")
	}

	vv := reflect.ValueOf(v)
	// If the stored value is a pointer and the target is its base type, dereference
	if vv.Kind() == reflect.Ptr && vv.Type().Elem().AssignableTo(rv.Elem().Type()) {
		rv.Elem().Set(vv.Elem())
		return nil
	}

	// Try direct set
	if vv.Type().AssignableTo(rv.Elem().Type()) {
		rv.Elem().Set(vv)
		return nil
	}

	return fmt.Errorf("cannot assign type %s to %s", vv.Type(), rv.Elem().Type())
}

func (s *RistrettoStore) Set(ctx context.Context, key string, val any, ttl time.Duration) error {
	// If val is a pointer, store the underlying value to avoid storing temporary pointers
	vv := reflect.ValueOf(val)
	if vv.Kind() == reflect.Ptr && !vv.IsNil() {
		s.cache.SetWithTTL(key, vv.Elem().Interface(), 1, ttl)
	} else {
		s.cache.SetWithTTL(key, val, 1, ttl)
	}
	return nil
}

func (s *RistrettoStore) Delete(ctx context.Context, key string) error {
	s.cache.Del(key)
	return nil
}

func (s *RistrettoStore) Flush(ctx context.Context) error {
	s.cache.Clear()
	return nil
}
