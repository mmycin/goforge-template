package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// AllConfig holds all configuration values
type AllConfig struct {
	App     AppConfig
	DB      DBConfig
	GRPC    GRPCConfig
	HTTP    HTTPConfig
	Encrypt EncryptConfig
	Log     LogConfig
}

// Package-level config variables (temporary for backwards compatibility)
var (
	App     AppConfig
	DB      DBConfig
	GRPC    GRPCConfig
	HTTP    HTTPConfig
	Encrypt EncryptConfig
	Log     LogConfig
)

// Load initializes and loads all configuration from .env file
func Load() (*AllConfig, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Warning: .env file not found, using defaults and environment variables")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var cfg AllConfig

	// Load all config sections
	var err error
	if cfg.App, err = loadAppConfig(); err != nil {
		return nil, err
	}
	if cfg.DB, err = loadDBConfig(); err != nil {
		return nil, err
	}
	if cfg.GRPC, err = loadGRPCConfig(); err != nil {
		return nil, err
	}
	if cfg.HTTP, err = loadHTTPConfig(); err != nil {
		return nil, err
	}
	if cfg.Encrypt, err = loadEncryptConfig(); err != nil {
		return nil, err
	}
	if cfg.Log, err = loadLogConfig(); err != nil {
		return nil, err
	}

	// Set package-level variables for backwards compatibility
	App = cfg.App
	DB = cfg.DB
	GRPC = cfg.GRPC
	HTTP = cfg.HTTP
	Encrypt = cfg.Encrypt
	Log = cfg.Log

	return &cfg, nil
}
