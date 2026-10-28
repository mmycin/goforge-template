// Package database provides the seeder engine for GoForge applications.
// Seeders implement the Seeder interface and are registered into a global
// registry so individual seeders can be executed by name via the CLI.
package database

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
)

// Seeder is the contract every database seeder must satisfy.
type Seeder interface {
	// Name returns the unique identifier for this seeder, e.g. "UserSeeder".
	Name() string
	// Run receives an open *gorm.DB and performs the seed logic.
	// It must be idempotent — running it multiple times should be safe.
	Run(db *gorm.DB) error
}

// ─── Global Seeder Registry ───────────────────────────────────────────────────

var (
	seederRegistryMu sync.RWMutex
	seederRegistry   = map[string]Seeder{}
)

// RegisterSeeder adds a Seeder to the global seeder registry.
// It is safe to call from multiple init() functions.
// Panics if two seeders share the same Name().
func RegisterSeeder(s Seeder) {
	seederRegistryMu.Lock()
	defer seederRegistryMu.Unlock()
	if _, ok := seederRegistry[s.Name()]; ok {
		panic(fmt.Sprintf("seeder: duplicate registration for %q", s.Name()))
	}
	seederRegistry[s.Name()] = s
}

// ResolveSeeder looks up a seeder by name in the global registry.
// Returns nil, false if the name is not registered.
func ResolveSeeder(name string) (Seeder, bool) {
	seederRegistryMu.RLock()
	defer seederRegistryMu.RUnlock()
	s, ok := seederRegistry[name]
	return s, ok
}

// AllSeederNames returns the names of all registered seeders.
func AllSeederNames() []string {
	seederRegistryMu.RLock()
	defer seederRegistryMu.RUnlock()
	names := make([]string, 0, len(seederRegistry))
	for n := range seederRegistry {
		names = append(names, n)
	}
	return names
}

// ─── Runner ──────────────────────────────────────────────────────────────────

// Call executes the provided seeders in the given order, wrapping them in a
// single database transaction. If any seeder returns an error the transaction
// is rolled back and the error is returned.
//
// Example usage inside DatabaseSeeder.Run:
//
//	return database.Call(db,
//	    &UserSeeder{},
//	    &TodoSeeder{},
//	)
func Call(db *gorm.DB, seeders ...Seeder) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for _, s := range seeders {
			start := time.Now()
			fmt.Printf("  → Running %-40s", s.Name()+"…")
			if err := s.Run(tx); err != nil {
				fmt.Println(" FAILED")
				return fmt.Errorf("seeder %q: %w", s.Name(), err)
			}
			fmt.Printf(" done (%dms)\n", time.Since(start).Milliseconds())
		}
		return nil
	})
}
