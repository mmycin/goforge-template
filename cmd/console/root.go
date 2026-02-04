package console

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goforge",
	Short: "GoForge - A powerful Go application framework",
	Long: `GoForge is a comprehensive Go application framework that provides
tools for database migrations, code generation, and service scaffolding.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Load environment variables before any command runs
		loadEnv()
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add all subcommands
	rootCmd.AddCommand(makeServiceCmd)
	rootCmd.AddCommand(genMigrationCmd)
	rootCmd.AddCommand(genSqlcCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(loaderCmd)
}
