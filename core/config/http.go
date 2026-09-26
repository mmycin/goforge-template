package config

import "github.com/spf13/viper"

// HTTPConfig holds HTTP server configuration
type HTTPConfig struct {
	RateLimitPerMinute int
	CORSAllowedOrigins string

	// TLS configuration
	TLSEnabled  bool
	TLSCertFile string
	TLSKeyFile  string
}

func loadHTTPConfig() (HTTPConfig, error) {
	return HTTPConfig{
		RateLimitPerMinute: viper.GetInt("RATE_LIMIT_PER_MINUTE"),
		CORSAllowedOrigins: viper.GetString("CORS_ALLOWED_ORIGINS"),

		TLSEnabled:  viper.GetBool("TLS_ENABLED"),
		TLSCertFile: viper.GetString("TLS_CERT_FILE"),
		TLSKeyFile:  viper.GetString("TLS_KEY_FILE"),
	}, nil
}
