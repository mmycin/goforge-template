package boot

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mmycin/goforge/boot/server"
	"github.com/mmycin/goforge/core/config"
	"github.com/mmycin/goforge/core/database"
	_ "github.com/mmycin/goforge/core/database/drivers" // registers all DB driver factories via init()
	"github.com/mmycin/goforge/core/events"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

// App holds all core application dependencies.
type App struct {
	Config   *config.AllConfig
	Database *database.Database
	EventBus *events.EventBus
	// Servers is the list of servers to run concurrently.
	// Add or remove a server here — kernel.go never needs to change again.
	Servers []server.Server
}

// NewApp creates and bootstraps a new App instance.
// It always starts the HTTP server and conditionally starts gRPC when enabled
// in config. Use NewAppWithServices for finer-grained control.
func NewApp(httpRouters []server.Router, grpcRegistries []server.GRPCRegistry) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	db, err := database.Connect(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	servers := buildServers(cfg, httpRouters, grpcRegistries)

	return &App{
		Config:   cfg,
		Database: db,
		EventBus: events.New(),
		Servers:  servers,
	}, nil
}

// NewAppWithServices creates a new App instance with pre-initialized dependencies.
// bus is created in core/console/serve.go and passed in so listeners registered
// during InitializeServices are already attached before Start() is called.
func NewAppWithServices(
	cfg  *config.AllConfig,
	db   *database.Database,
	bus  *events.EventBus,
	servers []server.Server,
) (*App, error) {
	return &App{
		Config:   cfg,
		Database: db,
		EventBus: bus,
		Servers:  servers,
	}, nil
}

// buildServers constructs the default server slice from config.
// HTTP is always included; gRPC is added when GRPC.Enable is true.
func buildServers(cfg *config.AllConfig, httpRouters []server.Router, grpcRegistries []server.GRPCRegistry) []server.Server {
	servers := []server.Server{
		server.NewHTTPServer(cfg, httpRouters),
	}
	if cfg.GRPC.Enable {
		servers = append(servers, server.NewGRPCServer(cfg, grpcRegistries))
	}
	return servers
}

// Start runs all registered servers concurrently and manages graceful shutdown.
// Every server is started in its own goroutine paired with a shutdown goroutine
// that fires when the OS signal context is cancelled.
// To add a new server type, implement server.Server and append it to App.Servers
// before calling Start — no changes to this method are needed.
func (app *App) Start() error {
	log.Info().Msg("Initializing GoForge Server")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, gCtx := errgroup.WithContext(ctx)

	for _, srv := range app.Servers {
		srv := srv // capture loop variable

		// Start the server.
		g.Go(func() error {
			if err := srv.Start(); err != nil {
				log.Error().Err(err).Str("server", srv.Name()).Msg("Server failed")
				return err
			}
			return nil
		})

		// Paired shutdown goroutine.
		g.Go(func() error {
			<-gCtx.Done()
			log.Info().Str("server", srv.Name()).Msg("Shutting down server...")

			shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			if err := srv.Stop(shutdownCtx); err != nil {
				log.Error().Err(err).Str("server", srv.Name()).Msg("Server forced to shutdown")
				return err
			}
			log.Info().Str("server", srv.Name()).Msg("Server stopped gracefully")
			return nil
		})
	}

	// Wait for all servers and their shutdown goroutines to return.
	if err := g.Wait(); err != nil {
		if err != context.Canceled {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			return err
		}
	}

	// Drain the event bus after all servers are down — no new events can arrive.
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
