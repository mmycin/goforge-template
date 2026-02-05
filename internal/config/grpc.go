package config

// GRPCConfig holds gRPC server configuration
type GRPCConfig struct {
	Enable     bool
	Reflection bool
	Host       string
	Port       int
}
