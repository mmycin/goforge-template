package seeders

import (
	"github.com/mmycin/goforge/core/database"
	"github.com/mmycin/goforge/internal/services/todo"
	"gorm.io/gorm"
)

// TodoSeeder seeds the initial todo items into the database.
// It is idempotent: it uses FirstOrCreate so re-running does not create duplicates.
type TodoSeeder struct{}

func (s *TodoSeeder) Name() string { return "TodoSeeder" }

func (s *TodoSeeder) Run(db *gorm.DB) error {
	items := []todo.Todo{
		{
			Title:       "Explore GoForge",
			Description: "Browse routes, events, services, and the console layer",
			Completed:   false,
		},
		{
			Title:       "Run migrations & seeders",
			Description: "Database schema and initial data are now in sync",
			Completed:   true,
		},
		{
			Title:       "Build something great",
			Description: "Start scaffolding your own services with goforge gen:service",
			Completed:   false,
		},
	}

	for i := range items {
		lookup := todo.Todo{Title: items[i].Title}
		if err := db.Where(lookup).FirstOrCreate(&items[i]).Error; err != nil {
			return err
		}
	}

	return nil
}

// Ensure TodoSeeder is registered at package init time.
func init() {
	database.RegisterSeeder(&TodoSeeder{})
}
