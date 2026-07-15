package config

import "github.com/spf13/viper"

// HTTPConfig holds HTTP server configuration
type HTTPConfig struct {
	RateLimitPerMinute int
	CORSAllowedOrigins string
}

func loadHTTPConfig() error {
	HTTP = HTTPConfig{
		RateLimitPerMinute: viper.GetInt("RATE_LIMIT_PER_MINUTE"),
		CORSAllowedOrigins: viper.GetString("CORS_ALLOWED_ORIGINS"),
	}
	return nil
}
