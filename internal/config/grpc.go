package config

import "github.com/spf13/viper"

// GRPCConfig holds gRPC server configuration
type GRPCConfig struct {
	Enable     bool
	Reflection bool
	Host       string
	Port       int
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
