package drivers

import (
	"github.com/mmycin/goforge/core/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	database.Register("postgres", func(dsn string) gorm.Dialector {
		return postgres.Open(dsn)
	})
}
