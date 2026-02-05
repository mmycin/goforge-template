package config

// CacheConfig holds cache configuration
type CacheConfig struct {
	Enabled  bool
	Driver   string
	TTL      string
	MaxItems int
	MaxCost  string
}
