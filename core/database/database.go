// Package database provides the Database wrapper and connection helpers
// for GoForge applications.
package database

import (
	"fmt"

	"gorm.io/gorm"
)

// Database holds the active database connections.
// Inject *Database into services the same way you would inject any other dependency.
type Database struct {
	Gorm *gorm.DB
	// Sqlc will be set when sqlc generated code is wired in.
	Sqlc any
}

// NewDatabase wraps an open *gorm.DB in a Database instance.
func NewDatabase(gormDB *gorm.DB) *Database {
	return &Database{Gorm: gormDB}
}

// GetGorm returns the underlying *gorm.DB with nil-safety checks.
// Always use this inside service methods rather than accessing .Gorm directly.
func (d *Database) GetGorm() (*gorm.DB, error) {
	if d == nil {
		return nil, fmt.Errorf("database: instance is nil")
	}
	if d.Gorm == nil {
		return nil, fmt.Errorf("database: gorm instance not initialized")
	}
	return d.Gorm, nil
}
