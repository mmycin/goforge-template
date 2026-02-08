package console

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mmycin/goforge/internal/config"
	"github.com/mmycin/goforge/internal/database"
	"github.com/mmycin/goforge/internal/server"
	"github.com/mmycin/goforge/internal/services"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server(s)",
	Long:  `Launch the REST API server and (optionally) the gRPC server with graceful shutdown support.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("Initializing GoForge Server")

		// Create context that listens for the interrupt signal from the OS.
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		// Error group for managing multiple server goroutines
		g, gCtx := errgroup.WithContext(ctx)

		log.Info().Msg("→ Connecting to database...")
		if err := database.Connect(); err != nil {
			log.Error().Err(err).Msg("Database connection failed")
			os.Exit(1)
		}

		// --- HTTP Server ---
		routers := services.GetRouters()
		httpSrv := server.NewHTTPServer(routers)

		g.Go(func() error {
			if err := httpSrv.Start(); err != nil && err != http.ErrServerClosed {
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

			if err := httpSrv.Stop(shutdownCtx); err != nil {
				log.Error().Err(err).Msg("HTTP server forced to shutdown")
				return err
			}
			log.Info().Msg("HTTP server stopped gracefully")
			return nil
		})

		// --- gRPC Server ---
		if config.GRPC.Enable {
			grpcSrv := server.NewGRPCServer()

			g.Go(func() error {
				if err := grpcSrv.Start(); err != nil {
					log.Error().Err(err).Msg("gRPC server failed")
					return err
				}
				return nil
			})

			g.Go(func() error {
				<-gCtx.Done()
				// gRPC GracefulStop waits for all RPCs to finish, so we do it in a goroutine
				// to allow the errgroup to proceed if needed, but here we block until done
				// or we could use a timeout wrapper if we wanted strictly bounded shutdown.
				// However, GracefulStop doesn't take a context.
				log.Info().Msg("Shutting down gRPC server...")

				// Standard gRPC GracefulStop block indefinitely if connections don't close.
				// For robustness, one might want to use Stop() after a timeout.
				// Here we trust GracefulStop for now.
				grpcSrv.Stop()
				log.Info().Msg("gRPC server stopped gracefully")
				return nil
			})
		}

		// Wait for all servers to stop
		if err := g.Wait(); err != nil {
			// If error is context.Canceled, it's a normal shutdown usually
			if err != context.Canceled {
				fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			}
		}
		log.Info().Msg("Server shutdown complete")
	},
}
