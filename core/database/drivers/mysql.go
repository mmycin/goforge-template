package drivers

import (
	"github.com/mmycin/goforge/core/database"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	database.Register("mysql", func(dsn string) gorm.Dialector {
		return mysql.Open(dsn)
	})
}
