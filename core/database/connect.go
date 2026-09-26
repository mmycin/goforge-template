package database

import (
	"fmt"

	"github.com/mmycin/goforge/core/config"
	"gorm.io/gorm"
)

// Connect opens a database connection from the given config and returns a
// ready-to-use *Database.
//
// The driver is selected from cfg.Connection. The matching driver file in
// core/database/drivers/ must be present — its init() registers the factory.
// If the driver file was deleted (e.g. by the CLI during project setup),
// Connect returns a descriptive error rather than panicking.
func Connect(cfg config.DBConfig) (*Database, error) {
	driver, dsn := buildDSN(cfg)

	factory, ok := drivers[driver]
	if !ok {
		return nil, fmt.Errorf(
			"database: no driver registered for %q — "+
				"check that core/database/drivers/%s.go exists and is compiled in",
			driver, driver,
		)
	}

	gormDB, err := gorm.Open(factory(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("database: open %s connection: %w", driver, err)
	}

	return NewDatabase(gormDB), nil
}
