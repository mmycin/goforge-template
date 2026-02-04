package console

import (
	"fmt"

	"github.com/spf13/cobra"
)

var genMigrationCmd = &cobra.Command{
	Use:   "gen:migration [name]",
	Short: "Create a new database migration",
	Long:  `Generate a new database migration file with the specified name.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		fmt.Printf("Creating migration: %s\n", name)
		makeMigration(name)
	},
}
