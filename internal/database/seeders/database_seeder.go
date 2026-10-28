package seeders

import (
	"github.com/mmycin/goforge/core/database"
	"gorm.io/gorm"
)

// DatabaseSeeder is the root seeder. It is executed by default when running
// `goforge db:seed` without a --class flag.
//
// Add child seeders here in the order they should run. database.Call wraps
// all of them in a single transaction — if any fails, the entire seed is
// rolled back.
//
// Example:
//
//	return database.Call(db,
//	    &UserSeeder{},
//	    &TodoSeeder{},
//	    &CategorySeeder{},
//	)
type DatabaseSeeder struct{}

func (s *DatabaseSeeder) Name() string { return "DatabaseSeeder" }

func (s *DatabaseSeeder) Run(db *gorm.DB) error {
	return database.Call(db,
		&TodoSeeder{},
		// Add more seeders here as your application grows.
	)
}

// Ensure DatabaseSeeder is registered at package init time.
func init() {
	database.RegisterSeeder(&DatabaseSeeder{})
}
