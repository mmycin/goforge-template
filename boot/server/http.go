package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mmycin/goforge/boot/server/middleware"
	"github.com/mmycin/goforge/core/config"
	"github.com/rs/zerolog/log"
)

// Server is the common interface every server type must satisfy.
// Implement this interface to plug a new server kind (HTTP, gRPC, GraphQL, …)
// into the application without touching boot/kernel.go.
type Server interface {
	// Name returns a human-readable identifier used in log messages.
	Name() string
	// Start begins accepting connections. It must block until the server
	// is stopped or encounters a fatal error.
	Start() error
	// Stop shuts the server down gracefully. The provided context sets the
	// deadline; implementations should force-stop when the context expires.
	Stop(ctx context.Context) error
}

// Router interface for service route registration
type Router interface {
	Register(engine *gin.Engine)
}

// HTTPServer represents the HTTP server
type HTTPServer struct {
	engine    *gin.Engine
	server    *http.Server
	tlsConfig *config.HTTPConfig
}

// NewHTTPServer creates a new HTTP server instance
func NewHTTPServer(cfg *config.AllConfig, routers []Router) *HTTPServer {
	// Set Gin mode based on debug config
	if !cfg.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin engine
	engine := gin.New()

	// Add middleware
	engine.Use(middleware.Recovery())
	engine.Use(middleware.CustomLogger())
	engine.Use(middleware.CORS())
	engine.Use(middleware.RateLimiter())
	engine.Use(middleware.ErrorHandler())

	// Health check endpoint (Public)
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": cfg.App.Name,
			"version": cfg.App.Version,
		})
	})

	// Apply AppKey middleware globally for all subsequent routes
	engine.Use(middleware.AppKey())

	// Register all service routers
	for _, router := range routers {
		router.Register(engine)
	}

	// Create HTTP server
	addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	// Apply TLS configuration when enabled
	if cfg.HTTP.TLSEnabled {
		srv.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	return &HTTPServer{
		engine:    engine,
		server:    srv,
		tlsConfig: &cfg.HTTP,
	}
}

// Name returns the server identifier used in log messages.
func (s *HTTPServer) Name() string {
	if s.tlsConfig.TLSEnabled {
		return "HTTPS"
	}
	return "HTTP"
}

// Start starts the HTTP server (plain HTTP or HTTPS depending on config)
func (s *HTTPServer) Start() error {
	if s.tlsConfig.TLSEnabled {
		if s.tlsConfig.TLSCertFile == "" || s.tlsConfig.TLSKeyFile == "" {
			return fmt.Errorf("TLS is enabled but TLS_CERT_FILE and/or TLS_KEY_FILE are not set")
		}
		log.Info().Str("addr", s.server.Addr).Msg("Starting HTTPS server")
		if err := s.server.ListenAndServeTLS(s.tlsConfig.TLSCertFile, s.tlsConfig.TLSKeyFile); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("HTTPS server ListenAndServeTLS failed: %w", err)
		}
		return nil
	}

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
