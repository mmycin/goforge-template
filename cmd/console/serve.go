package console

import (
	"fmt"
	"os"

	"github.com/mmycin/goforge/cmd/api"
	"github.com/mmycin/goforge/internal/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP API server",
	Long:  `Launch the Gin-powered REST API server with graceful shutdown support.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("Initializing GoForge API Server")

		// Get all routers
		routers := api.GetRouters()

		// Create and start server
		srv := server.NewHTTPServer(routers)
		if err := srv.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}
