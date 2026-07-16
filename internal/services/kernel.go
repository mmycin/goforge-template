package services

import (
	"github.com/mmycin/goforge/boot/server"
	"github.com/mmycin/goforge/core/database"
	"github.com/mmycin/goforge/core/events"
	"github.com/mmycin/goforge/internal/services/todo"
)

// ServicesConfig contains all dependencies needed to initialize services.
type ServicesConfig struct {
	DB  *database.Database
	Bus *events.EventBus
}

// InitializeServices initializes all services, registers event listeners,
// and returns HTTP routers and gRPC registries.
func InitializeServices(cfg ServicesConfig) ([]server.Router, []server.GRPCRegistry) {
	// --- Todo service ---
	todoService := todo.NewTodoService(cfg.DB, cfg.Bus)
	todoHandler := todo.NewTodoHandler(todoService)
	todoRoutes := todo.NewTodoRoutes(todoHandler)
	todoDocs := todo.NewTodoDocs()

	// --- Event listeners ---
	// Each listener is registered against its topic constant.
	// Add new listeners here as the application grows.
	events.On(cfg.Bus, todo.EventTodoCreated, &todo.AuditOnCreated{DB: cfg.DB})
	events.On(cfg.Bus, todo.EventTodoUpdated, &todo.AuditOnUpdated{DB: cfg.DB})
	events.On(cfg.Bus, todo.EventTodoDeleted, &todo.AuditOnDeleted{DB: cfg.DB})

	routers := []server.Router{
		todoRoutes,
		todoDocs,
	}

	grpcRegistries := []server.GRPCRegistry{}

	return routers, grpcRegistries
}

// Model returns all GORM models to be registered for migrations.
func Model() []any {
	return []any{
		&todo.Todo{},
		&todo.AuditLog{},
	}
}
