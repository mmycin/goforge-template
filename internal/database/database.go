package database

import (
	"fmt"
	"strings"

	"github.com/mmycin/goforge/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type Database struct {
	Gorm *gorm.DB
	// Sqlc field will be added when generated code is available
	Sqlc any
}

// NewDatabase creates a new Database instance
func NewDatabase(gormDB *gorm.DB) *Database {
	return &Database{Gorm: gormDB}
}

// GetGorm returns the Gorm instance with nil check
func (d *Database) GetGorm() (*gorm.DB, error) {
	if d == nil {
		return nil, fmt.Errorf("database instance is nil")
	}
	if d.Gorm == nil {
		return nil, fmt.Errorf("database gorm instance not initialized")
	}
	return d.Gorm, nil
}

// Connect takes DBConfig and returns *Database
func Connect(cfg config.DBConfig) (*Database, error) {
	dbName := cfg.Name
	dbDriver := cfg.Connection
	dbDsn := ""

	if dbDriver == "" {
		if strings.HasSuffix(dbName, ".db") {
			dbDriver = "sqlite"
			dbDsn = dbName
		} else {
			dbDriver = "sqlite"
		}
	}

	if dbDsn == "" {
		switch dbDriver {
		case "mysql":
			dbDsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
		case "postgres":
			dbDsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
				cfg.Host, cfg.Username, cfg.Password, cfg.Name, cfg.Port)
		case "sqlite":
			dbDsn = dbName
		case "sqlserver":
			dbDsn = fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
				cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
		}
	}

	var dialector gorm.Dialector
	switch dbDriver {
	case "mysql":
		dialector = mysql.Open(dbDsn)
	case "postgres":
		dialector = postgres.Open(dbDsn)
	case "sqlite":
		dialector = sqlite.Open(dbDsn)
	case "sqlserver":
		dialector = sqlserver.Open(dbDsn)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", dbDriver)
	}

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return NewDatabase(gormDB), nil
}
