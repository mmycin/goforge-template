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

var DB *Database

type Database struct {
	Gorm *gorm.DB
}

func Connect() error {
	dbName := config.DB.Name
	dbDriver := config.DB.Connection
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
				config.DB.Username, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.Name)
		case "postgres":
			dbDsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
				config.DB.Host, config.DB.Username, config.DB.Password, config.DB.Name, config.DB.Port)
		case "sqlite":
			dbDsn = dbName
		case "sqlserver":
			dbDsn = fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
				config.DB.Username, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.Name)
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
		return fmt.Errorf("unsupported database driver: %s", dbDriver)
	}

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return err
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = &Database{
		Gorm: gormDB,
	}

	return nil
}
