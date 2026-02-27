package config

import "github.com/spf13/viper"

// DBConfig holds database configuration
type DBConfig struct {
	Connection string
	Name       string
	DevName    string
	Host       string
	Port       int
	Username   string
	Password   string
	Migrator   string
}

func loadDBConfig() error {
	DB = DBConfig{
		Connection: viper.GetString("DB_CONNECTION"),
		Name:       viper.GetString("DB_NAME"),
		DevName:    viper.GetString("DB_DEV_NAME"),
		Host:       viper.GetString("DB_HOST"),
		Port:       viper.GetInt("DB_PORT"),
		Username:   viper.GetString("DB_USERNAME"),
		Password:   viper.GetString("DB_PASSWORD"),
		Migrator:   viper.GetString("DB_MIGRATOR"),
	}
	return nil
}
