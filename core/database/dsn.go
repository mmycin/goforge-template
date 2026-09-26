package database

import (
	"fmt"

	"github.com/mmycin/goforge/core/config"
)

// buildDSN returns the driver name and DSN string for the given config.
// The driver name comes directly from cfg.Connection — no silent fallback.
// If cfg.Connection is blank, the caller gets an empty driver name which
// Connect will reject with a clear error.
func buildDSN(cfg config.DBConfig) (driver, dsn string) {
	driver = cfg.Connection
	name := cfg.Name

	switch driver {
	case "mysql":
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, name,
		)
	case "postgres":
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
			cfg.Host, cfg.Username, cfg.Password, name, cfg.Port,
		)
	case "sqlite":
		dsn = name
	case "sqlserver":
		dsn = fmt.Sprintf(
			"sqlserver://%s:%s@%s:%d?database=%s",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, name,
		)
	}

	return driver, dsn
}
