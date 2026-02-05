package server

import (
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	routers []Router
	mu      sync.Mutex
)

// Register adds a router to the global registry
func Register(r Router) {
	mu.Lock()
	defer mu.Unlock()
	routers = append(routers, r)
}

// GetRegisteredRouters returns all registered routers
func GetRegisteredRouters() []Router {
	mu.Lock()
	defer mu.Unlock()

	// Return a copy to avoid race conditions
	r := make([]Router, len(routers))
	copy(r, routers)
	return r
}

// RegisterAll registers all routers in the engine
func RegisterAll(engine *gin.Engine) {
	for _, r := range GetRegisteredRouters() {
		r.Register(engine)
	}
}
