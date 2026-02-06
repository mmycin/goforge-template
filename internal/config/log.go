package config

import "github.com/spf13/viper"

// LogConfig holds logging configuration
type LogConfig struct {
	Type   string
	Level  string
	Path   string
	Format string
}

func loadLogConfig() error {
	Log = LogConfig{
		Type:   viper.GetString("LOG_TYPE"),
		Level:  viper.GetString("LOG_LEVEL"),
		Path:   viper.GetString("LOG_PATH"),
		Format: viper.GetString("LOG_FORMAT"),
	}
	return nil
}
