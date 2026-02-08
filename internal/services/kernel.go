package services

import (
	"github.com/mmycin/goforge/internal/server"
	"github.com/mmycin/goforge/internal/services/todo"
)

// GetRouters returns all service routers to be registered
func GetRouters() []server.Router {
	return server.GetRegisteredRouters()
}

// Model returns all models to be registered with GORM
func Model() []any {
	return []any{
		&todo.Todo{},
	}
}
