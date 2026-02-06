package server

import (
	"context"
	"fmt"
	"net/http"
	"sync" // Added sync for registry

	"github.com/gin-gonic/gin"
	"github.com/mmycin/goforge/internal/config"
	"github.com/mmycin/goforge/internal/server/middleware"
	"github.com/rs/zerolog/log"
)

// Router interface for service route registration
type Router interface {
	Register(engine *gin.Engine)
}

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

	r := make([]Router, len(routers))
	copy(r, routers)
	return r
}

// HTTPServer represents the HTTP server
type HTTPServer struct {
	engine *gin.Engine
	server *http.Server
}

// NewHTTPServer creates a new HTTP server instance
func NewHTTPServer(routers []Router) *HTTPServer {
	// Set Gin mode based on debug config
	if !config.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin engine
	engine := gin.New()

	// Add middleware
	engine.Use(middleware.Recovery())
	engine.Use(middleware.CustomLogger())
	engine.Use(middleware.CORS())
	engine.Use(middleware.RateLimiter())

	// Health check endpoint (Public)
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": config.App.Name,
			"version": config.App.Version,
		})
	})

	// Apply AppKey middleware globally for all subsequent routes
	engine.Use(middleware.AppKey())

	// Register all service routers
	for _, router := range routers {
		router.Register(engine)
	}

	// Create HTTP server
	addr := fmt.Sprintf("%s:%d", config.App.Host, config.App.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	return &HTTPServer{
		engine: engine,
		server: server,
	}
}

// Start starts the HTTP server
func (s *HTTPServer) Start() error {
	log.Info().Str("addr", s.server.Addr).Msg("Starting HTTP server")
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server ListenAndServe failed: %w", err)
	}
	return nil
}

// Stop stops the HTTP server
func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// Engine returns the underlying Gin engine
func (s *HTTPServer) Engine() *gin.Engine {
	return s.engine
}
