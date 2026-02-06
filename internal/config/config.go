package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

var (
	App     AppConfig
	DB      DBConfig
	GRPC    GRPCConfig
	HTTP    HTTPConfig
	Cache   CacheConfig
	Encrypt EncryptConfig
	Log     LogConfig
)

// Load initializes and loads all configuration from .env file
func Load() error {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Warning: .env file not found, using defaults and environment variables")
		} else {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Load all config sections
	if err := loadAppConfig(); err != nil {
		return err
	}
	if err := loadDBConfig(); err != nil {
		return err
	}
	if err := loadGRPCConfig(); err != nil {
		return err
	}
	if err := loadHTTPConfig(); err != nil {
		return err
	}
	if err := loadCacheConfig(); err != nil {
		return err
	}
	if err := loadEncryptConfig(); err != nil {
		return err
	}
	if err := loadLogConfig(); err != nil {
		return err
	}

	return nil
}
