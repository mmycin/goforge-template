package config

// EncryptConfig holds encryption configuration
type EncryptConfig struct {
	Key    string
	Rounds int
	Salt   string
}
