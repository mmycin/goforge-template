package services

import (
	"github.com/mmycin/goforge/internal/server"
	_ "github.com/mmycin/goforge/internal/services/todo"
)

// GetRouters returns all service routers to be registered
func GetRouters() []server.Router {
	return server.GetRegisteredRouters()
}
