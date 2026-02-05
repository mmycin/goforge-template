package config

// AppConfig holds application-level configuration
type AppConfig struct {
	Name    string
	Version string
	Debug   bool
	Key     string
	Host    string
	Port    int
}
