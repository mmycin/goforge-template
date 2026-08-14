package console

import (
	"fmt"
	"os"

	"github.com/mmycin/goforge/core/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goforge",
	Short: "GoForge - A powerful Go application framework",
	Long: `GoForge is a comprehensive Go application framework that provides
tools for database migrations, code generation, and service scaffolding.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Load configuration before any command runs
		if _, err := config.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
			os.Exit(1)
		}
		// Initialize logger after config is loaded
		initLogger()
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add all subcommands
	rootCmd.AddCommand(serveCmd)

	// Default command: running the binary with no args starts the server.
	rootCmd.RunE = serveCmd.RunE
}
