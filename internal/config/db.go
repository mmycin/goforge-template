package config

// DBConfig holds database configuration
type DBConfig struct {
	Connection string
	Name       string
	DevName    string
	Host       string
	Port       int
	Username   string
	Password   string
}
