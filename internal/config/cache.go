package config

import "github.com/spf13/viper"

// CacheConfig holds cache configuration
type CacheConfig struct {
	Enabled  bool
	Driver   string
	TTL      string
	MaxItems int
	MaxCost  string
	Redis    RedisConfig
}

// RedisConfig holds redis connection settings
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	Database int
}

func loadCacheConfig() error {
	Cache = CacheConfig{
		Enabled:  viper.GetBool("CACHE_ENABLED"),
		Driver:   viper.GetString("CACHE_DRIVER"),
		TTL:      viper.GetString("CACHE_TTL"),
		MaxItems: viper.GetInt("CACHE_MAX_ITEMS"),
		MaxCost:  viper.GetString("CACHE_MAX_COST"),
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetInt("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			Database: viper.GetInt("REDIS_DB"),
		},
	}
	return nil
}
