package console

import (
	"fmt"
	"os"

	"github.com/mmycin/goforge/internal/cache"
	"github.com/mmycin/goforge/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goforge",
	Short: "GoForge - A powerful Go application framework",
	Long: `GoForge is a comprehensive Go application framework that provides
tools for database migrations, code generation, and service scaffolding.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Load configuration before any command runs
		if err := config.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
			os.Exit(1)
		}
		// Initialize logger after config is loaded
		initLogger()

		// Initialize cache
		if err := cache.Connect(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Cache initialization failed: %v\n", err)
		}
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add all subcommands
	rootCmd.AddCommand(serveCmd)
}
