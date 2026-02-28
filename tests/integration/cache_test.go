package integration_tests

import (
	"context"
	"testing"
	"time"

	"github.com/mmycin/goforge/internal/cache"
)

func TestRistrettoStore(t *testing.T) {
	ctx := context.Background()
	store, err := cache.NewRistrettoStore(100, 1000)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	t.Run("SetAndGet", func(t *testing.T) {
		key := "test_key"
		val := "test_value"
		err := store.Set(ctx, key, val, 1*time.Minute)
		if err != nil {
			t.Errorf("Set failed: %v", err)
		}

		// Ristretto is eventually consistent, wait a bit
		time.Sleep(20 * time.Millisecond)

		var result string
		err = store.Get(ctx, key, &result)
		if err != nil {
			t.Errorf("Get failed: %v", err)
		}
		if result != val {
			t.Errorf("Expected %s, got %s", val, result)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		key := "del_key"
		_ = store.Set(ctx, key, "val", 1*time.Minute)
		time.Sleep(20 * time.Millisecond)

		_ = store.Delete(ctx, key)
		var res string
		err := store.Get(ctx, key, &res)
		if err == nil {
			t.Error("Expected error for deleted key, got nil")
		}
	})

	t.Run("Flush", func(t *testing.T) {
		_ = store.Set(ctx, "k1", "v1", 1*time.Minute)
		_ = store.Set(ctx, "k2", "v2", 1*time.Minute)
		time.Sleep(20 * time.Millisecond)

		_ = store.Flush(ctx)
		var res string
		if err := store.Get(ctx, "k1", &res); err == nil {
			t.Error("Expected error after flush")
		}
	})
}

func TestMultiStore(t *testing.T) {
	ctx := context.Background()
	l1, _ := cache.NewRistrettoStore(100, 1000)
	l2, _ := cache.NewRistrettoStore(100, 1000)
	multi := cache.NewMultiStore(l1, l2)

	t.Run("TieredGet", func(t *testing.T) {
		key := "tiered_key"
		val := "tiered_val"

		// Set only in L2
		_ = l2.Set(ctx, key, val, 1*time.Minute)
		time.Sleep(20 * time.Millisecond)

		// Get from Multi (should trigger L1 update)
		var result string
		err := multi.Get(ctx, key, &result)
		if err != nil {
			t.Fatalf("Multi Get failed: %v", err)
		}

		// Verify L1 was updated
		time.Sleep(20 * time.Millisecond)
		var l1Res string
		if err := l1.Get(ctx, key, &l1Res); err != nil {
			t.Errorf("L1 should have been updated: %v", err)
		}
	})

	t.Run("MultiSet", func(t *testing.T) {
		key := "set_key"
		val := "set_val"
		_ = multi.Set(ctx, key, val, 1*time.Minute)
		time.Sleep(20 * time.Millisecond)

		var r1, r2 string
		_ = l1.Get(ctx, key, &r1)
		_ = l2.Get(ctx, key, &r2)

		if r1 != val || r2 != val {
			t.Errorf("Expected both to have %s, got %s and %s", val, r1, r2)
		}
	})
}
