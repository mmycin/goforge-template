package services

import (
	"github.com/mmycin/goforge/boot/server"
	"github.com/mmycin/goforge/core/database"
	"github.com/mmycin/goforge/internal/services/todo"
)

// ServicesConfig contains all dependencies needed to initialize services
type ServicesConfig struct {
	DB *database.Database
}

// InitializeServices initializes all services and returns routers and gRPC registries
func InitializeServices(cfg ServicesConfig) ([]server.Router, []server.GRPCRegistry) {
	// Initialize todo service
	todoService := todo.NewTodoService(cfg.DB)
	todoHandler := todo.NewTodoHandler(todoService)
	todoRoutes := todo.NewTodoRoutes(todoHandler)
	todoDocs := todo.NewTodoDocs()

	// Collect routers
	routers := []server.Router{
		todoRoutes,
		todoDocs,
	}

	// Collect gRPC registries (empty for now)
	grpcRegistries := []server.GRPCRegistry{}

	return routers, grpcRegistries
}

// Model returns all models to be registered with GORM
func Model() []any {
	return []any{
		&todo.Todo{},
	}
}
