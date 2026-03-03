package console

import (
	"fmt"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/mmycin/goforge/internal/config"
	"github.com/mmycin/goforge/internal/database"
	"github.com/mmycin/goforge/internal/services"
	"github.com/spf13/cobra"
)

var loaderCmd = &cobra.Command{
	Use:   "loader",
	Short: "Run GORM schema loader",
	Long:  `Load and display GORM schema definitions.`,
	Run: func(cmd *cobra.Command, args []string) {
		runLoader()
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Execute GORM AutoMigrate to synchronize database schema with models.`,
	Run: func(cmd *cobra.Command, args []string) {
		runMigrate()
	},
}

func init() {
	rootCmd.AddCommand(loaderCmd)
	rootCmd.AddCommand(migrateCmd)
}

func runLoader() {
	models := services.Model()
	driver := config.DB.Connection
	if driver == "" {
		driver = "sqlite"
	}
	loader := gormschema.New(driver)
	stmts, err := loader.Load(models...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to load GORM schema: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, stmts)
}

func runMigrate() {
	fmt.Println("→ Starting GORM AutoMigrate...")
	if err := database.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Database connection failed: %v\n", err)
		os.Exit(1)
	}

	models := services.Model()
	if err := database.DB.Gorm.AutoMigrate(models...); err != nil {
		fmt.Fprintf(os.Stderr, "Error: AutoMigrate failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Database migration completed successfully")
}
