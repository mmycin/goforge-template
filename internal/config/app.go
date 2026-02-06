package config

import "github.com/spf13/viper"

// AppConfig holds application-level configuration
type AppConfig struct {
	Name    string
	Version string
	Debug   bool
	Key     string
	Host    string
	Port    int
}

func loadAppConfig() error {
	App = AppConfig{
		Name:    viper.GetString("APP_NAME"),
		Version: viper.GetString("APP_VERSION"),
		Debug:   viper.GetBool("APP_DEBUG"),
		Key:     viper.GetString("APP_KEY"),
		Host:    viper.GetString("APP_HOST"),
		Port:    viper.GetInt("APP_PORT"),
	}
	return nil
}
