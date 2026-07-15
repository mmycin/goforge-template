package config

import "github.com/spf13/viper"

// EncryptConfig holds encryption and hashing configuration
type EncryptConfig struct {
	Key    string
	Rounds int
	Salt   string
}

func loadEncryptConfig() (EncryptConfig, error) {
	return EncryptConfig{
		Key:    viper.GetString("ENCRYPTION_KEY"),
		Rounds: viper.GetInt("ENCRYPTION_ROUNDS"),
		Salt:   viper.GetString("ENCRYPTION_SALT"),
	}, nil
}
