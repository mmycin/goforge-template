package console

import (
	"fmt"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Execute GORM AutoMigrate to synchronize database schema with models.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running database migration...")
		migrateDB()
	},
}
