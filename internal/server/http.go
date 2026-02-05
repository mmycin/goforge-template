package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mmycin/goforge/internal/config"
	"github.com/mmycin/goforge/internal/server/middleware"
)

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

	// Health check endpoint (Public)
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": config.App.Name,
			"version": config.App.Version,
		})
	})

	// Authenticated routes
	auth := engine.Group("/")
	auth.Use(middleware.AppKey())
	{
		// Register all service routers to the authenticated group
		for _, router := range routers {
			router.Register(auth)
		}
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

// Start starts the HTTP server with graceful shutdown
func (s *HTTPServer) Start() error {
	// Channel to listen for errors from the server
	serverErrors := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		fmt.Printf("→ Starting HTTP server on %s\n", s.server.Addr)
		serverErrors <- s.server.ListenAndServe()
	}()

	// Channel to listen for interrupt signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal or an error
	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("server error: %w", err)
		}

	case sig := <-shutdown:
		fmt.Printf("\n→ Received signal: %v\n", sig)
		fmt.Println("→ Starting graceful shutdown...")

		// Create context with timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := s.server.Shutdown(ctx); err != nil {
			// Force close if graceful shutdown fails
			s.server.Close()
			return fmt.Errorf("server shutdown error: %w", err)
		}

		fmt.Println("✓ Server stopped gracefully")
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
