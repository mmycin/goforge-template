package boot

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mmycin/goforge/boot/server"
	"github.com/mmycin/goforge/internal/config"
	"github.com/mmycin/goforge/internal/database"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

// App holds all core application dependencies
type App struct {
	Config     *config.AllConfig
	Database   *database.Database
	HTTPServer *server.HTTPServer
	GRPCServer *server.GRPCServer
}

// NewApp creates and bootstraps a new App instance
func NewApp(httpRouters []server.Router, grpcRegistries []server.GRPCRegistry) (*App, error) {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Connect to database
	db, err := database.Connect(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize servers
	httpSrv := server.NewHTTPServer(cfg, httpRouters)
	grpcSrv := server.NewGRPCServer(cfg, grpcRegistries)

	return &App{
		Config:     cfg,
		Database:   db,
		HTTPServer: httpSrv,
		GRPCServer: grpcSrv,
	}, nil
}

// NewAppWithServices creates a new App instance with pre-initialized config and database
func NewAppWithServices(cfg *config.AllConfig, db *database.Database, httpRouters []server.Router, grpcRegistries []server.GRPCRegistry) (*App, error) {
	// Initialize servers
	httpSrv := server.NewHTTPServer(cfg, httpRouters)
	grpcSrv := server.NewGRPCServer(cfg, grpcRegistries)

	return &App{
		Config:     cfg,
		Database:   db,
		HTTPServer: httpSrv,
		GRPCServer: grpcSrv,
	}, nil
}

// Start starts all servers and manages graceful shutdown
func (app *App) Start() error {
	log.Info().Msg("Initializing GoForge Server")

	// Create context that listens for interrupt signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Error group for managing server goroutines
	g, gCtx := errgroup.WithContext(ctx)

	// --- HTTP Server ---
	g.Go(func() error {
		if err := app.HTTPServer.Start(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("HTTP server failed")
			return err
		}
		return nil
	})

	g.Go(func() error {
		<-gCtx.Done()
		log.Info().Msg("Shutting down HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := app.HTTPServer.Stop(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("HTTP server forced to shutdown")
			return err
		}
		log.Info().Msg("HTTP server stopped gracefully")
		return nil
	})

	// --- gRPC Server ---
	if app.Config.GRPC.Enable {
		g.Go(func() error {
			if err := app.GRPCServer.Start(); err != nil {
				log.Error().Err(err).Msg("gRPC server failed")
				return err
			}
			return nil
		})

		g.Go(func() error {
			<-gCtx.Done()
			log.Info().Msg("Shutting down gRPC server...")
			app.GRPCServer.StopWithTimeout(10 * time.Second)
			return nil
		})
	}

	// Wait for all servers to stop
	if err := g.Wait(); err != nil {
		if err != context.Canceled {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			return err
		}
	}
	log.Info().Msg("Server shutdown complete")
	return nil
}
