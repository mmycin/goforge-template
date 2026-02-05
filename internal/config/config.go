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

func loadDBConfig() error {
	DB = DBConfig{
		Connection: viper.GetString("DB_CONNECTION"),
		Name:       viper.GetString("DB_NAME"),
		Host:       viper.GetString("DB_HOST"),
		Port:       viper.GetInt("DB_PORT"),
		Username:   viper.GetString("DB_USERNAME"),
		Password:   viper.GetString("DB_PASSWORD"),
	}
	return nil
}

func loadGRPCConfig() error {
	GRPC = GRPCConfig{
		Enable:     viper.GetBool("GRPC_ENABLE"),
		Reflection: viper.GetBool("GRPC_REFLECTION"),
		Host:       viper.GetString("GRPC_HOST"),
		Port:       viper.GetInt("GRPC_PORT"),
	}
	return nil
}

func loadHTTPConfig() error {
	HTTP = HTTPConfig{
		RateLimitPerMinute: viper.GetInt("RATE_LIMIT_PER_MINUTE"),
	}
	return nil
}

func loadCacheConfig() error {
	Cache = CacheConfig{
		Enabled:  viper.GetBool("CACHE_ENABLED"),
		Driver:   viper.GetString("CACHE_DRIVER"),
		TTL:      viper.GetString("CACHE_TTL"),
		MaxItems: viper.GetInt("CACHE_MAX_ITEMS"),
		MaxCost:  viper.GetString("CACHE_MAX_COST"),
	}
	return nil
}

func loadEncryptConfig() error {
	Encrypt = EncryptConfig{
		Key:    viper.GetString("ENCRYPTION_KEY"),
		Rounds: viper.GetInt("ENCRPTION_ROUNDS"),
		Salt:   viper.GetString("ENCRYPTION_SALT"),
	}
	return nil
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
