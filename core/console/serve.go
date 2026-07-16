package console

import (
	"os"

	"github.com/mmycin/goforge/boot"
	"github.com/mmycin/goforge/core/config"
	"github.com/mmycin/goforge/core/database"
	"github.com/mmycin/goforge/core/events"
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

		// Create the event bus before services so listeners can be registered
		// inside InitializeServices.
		bus := events.New()

		routers, grpcRegistries := services.InitializeServices(services.ServicesConfig{
			DB:  db,
			Bus: bus,
		})

		app, err := boot.NewAppWithServices(cfg, db, bus, routers, grpcRegistries)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create app")
			os.Exit(1)
		}

		if err := app.Start(); err != nil {
			log.Error().Err(err).Msg("Server failed")
			os.Exit(1)
		}
	},
}
