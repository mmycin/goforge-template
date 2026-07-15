package console

import (
	"os"

	"github.com/mmycin/goforge/boot"
	"github.com/mmycin/goforge/internal/config"
	"github.com/mmycin/goforge/internal/database"
	"github.com/mmycin/goforge/internal/services"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server(s)",
	Long:  `Launch the REST API server and (optionally) the gRPC server with graceful shutdown support.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("Initializing GoForge Server")

		// Load config and connect to DB via boot.NewApp
		// First, we need to create a temporary config to pass to services.InitializeServices?
		// Wait no, let's first load config, then connect to DB, then init services, then create app
		// Let's adjust boot.NewApp to take services? Or let's load config, connect DB, init services, then create app
		// Let's modify the approach here:
		cfg, err := config.Load()
		if err != nil {
			log.Error().Err(err).Msg("Failed to load config")
			os.Exit(1)
		}

		db, err := database.Connect(cfg.DB)
		if err != nil {
			log.Error().Err(err).Msg("Database connection failed")
			os.Exit(1)
		}

		// Initialize services
		routers, grpcRegistries := services.InitializeServices(services.ServicesConfig{
			DB: db,
		})

		// Create app
		app, err := boot.NewAppWithServices(cfg, db, routers, grpcRegistries)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create app")
			os.Exit(1)
		}

		// Start app
		if err := app.Start(); err != nil {
			log.Error().Err(err).Msg("Server failed")
			os.Exit(1)
		}
	},
}
