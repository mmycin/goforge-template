package cache

import (
	"fmt"

	"github.com/mmycin/goforge/internal/config"
)

// Connect initializes the global cache based on configuration
func Connect() error {
	var err error
	if !config.Cache.Enabled {
		return nil
	}

	driver := config.Cache.Driver
	fmt.Printf("→ Initializing cache with %s driver...\n", driver)

	switch driver {
	case "memory":
		Memory, err = NewRistrettoStore(int64(config.Cache.MaxItems), 1<<30)
		if err != nil {
			return err
		}
		Global = Memory
	case "redis":
		Redis = NewRedisStore(
			config.Cache.Redis.Host,
			config.Cache.Redis.Port,
			config.Cache.Redis.Password,
			config.Cache.Redis.Database,
		)
		Global = Redis
	case "both":
		Memory, err = NewRistrettoStore(int64(config.Cache.MaxItems), 1<<30)
		if err != nil {
			return err
		}
		Redis = NewRedisStore(
			config.Cache.Redis.Host,
			config.Cache.Redis.Port,
			config.Cache.Redis.Password,
			config.Cache.Redis.Database,
		)
		Global = NewMultiStore(Memory, Redis)
	default:
		return fmt.Errorf("unsupported cache driver: %s", driver)
	}

	fmt.Println("✓ Cache initialized successfully")
	return nil
}
