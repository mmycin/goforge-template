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
	"github.com/mmycin/goforge/core/config"
	"github.com/mmycin/goforge/core/database"
	"github.com/mmycin/goforge/core/events"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

// App holds all core application dependencies.
type App struct {
	Config     *config.AllConfig
	Database   *database.Database
	EventBus   *events.EventBus
	HTTPServer *server.HTTPServer
	GRPCServer *server.GRPCServer
}

// NewApp creates and bootstraps a new App instance without pre-initialized services.
func NewApp(httpRouters []server.Router, grpcRegistries []server.GRPCRegistry) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	db, err := database.Connect(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	httpSrv := server.NewHTTPServer(cfg, httpRouters)
	grpcSrv := server.NewGRPCServer(cfg, grpcRegistries)

	return &App{
		Config:     cfg,
		Database:   db,
		EventBus:   events.New(),
		HTTPServer: httpSrv,
		GRPCServer: grpcSrv,
	}, nil
}

// NewAppWithServices creates a new App instance with pre-initialized dependencies.
// bus is created in core/console/serve.go and passed in so listeners registered
// during InitializeServices are already attached before Start() is called.
func NewAppWithServices(
	cfg            *config.AllConfig,
	db             *database.Database,
	bus            *events.EventBus,
	httpRouters    []server.Router,
	grpcRegistries []server.GRPCRegistry,
) (*App, error) {
	httpSrv := server.NewHTTPServer(cfg, httpRouters)
	grpcSrv := server.NewGRPCServer(cfg, grpcRegistries)

	return &App{
		Config:     cfg,
		Database:   db,
		EventBus:   bus,
		HTTPServer: httpSrv,
		GRPCServer: grpcSrv,
	}, nil
}

// Start starts all servers and manages graceful shutdown.
func (app *App) Start() error {
	log.Info().Msg("Initializing GoForge Server")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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

	// Wait for all servers to stop.
	if err := g.Wait(); err != nil {
		if err != context.Canceled {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			return err
		}
	}

	// Drain the event bus after servers are down — no new events can arrive.
	log.Info().Msg("Draining event bus...")
	drainCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := app.EventBus.Shutdown(drainCtx); err != nil {
		log.Warn().Err(err).Msg("Event bus drain timed out; some events may be unprocessed")
	} else {
		log.Info().Msg("Event bus drained")
	}

	log.Info().Msg("Server shutdown complete")
	return nil
}
