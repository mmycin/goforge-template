package drivers

import (
	"github.com/mmycin/goforge/core/database"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func init() {
	database.Register("sqlserver", func(dsn string) gorm.Dialector {
		return sqlserver.Open(dsn)
	})
}
