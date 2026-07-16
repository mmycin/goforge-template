package database

import "gorm.io/gorm"

// DriverFactory creates a gorm.Dialector for a given DSN.
// Each driver file in drivers/ registers one factory via init().
type DriverFactory func(dsn string) gorm.Dialector

// drivers is the registry of available database drivers.
// Populated at program start by the init() functions in core/database/drivers/.
// Accessing this map is safe after init() runs and before any concurrent reads.
var drivers = map[string]DriverFactory{}

// Register adds a driver factory to the registry.
// Called exclusively from driver init() functions — never call this directly.
func Register(name string, factory DriverFactory) {
	drivers[name] = factory
}
