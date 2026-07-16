package drivers

import (
	"github.com/mmycin/goforge/core/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	database.Register("sqlite", func(dsn string) gorm.Dialector {
		return sqlite.Open(dsn)
	})
}
