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
	"github.com/rs/zerolog/log"
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
	fmt.Println("→ Starting server initialization...")

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Info().Str("addr", s.server.Addr).Msg("Starting HTTP server")
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("HTTP server ListenAndServe failed")
		}
	}()

	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	fmt.Println("→ Server is running. Press Ctrl+C to stop.")

	// Block until we receive a signal
	sig := <-quit
	log.Info().Str("signal", sig.String()).Msg("Shutdown signal received")
	fmt.Printf("\n→ Received signal: %v. Initiating graceful shutdown...\n", sig)

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := s.server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
		fmt.Printf("! Server forced to shutdown: %v\n", err)
		return err
	}

	log.Info().Msg("Server stopped gracefully")
	fmt.Println("✓ Server stopped gracefully")

	// Small delay to ensure logs are written
	time.Sleep(500 * time.Millisecond)

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
