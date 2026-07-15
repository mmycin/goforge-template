package console

import (
	"fmt"
	"os"
	"os/exec"

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
	Long:  `Execute migrations to synchronize database schema with models.`,
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
	// Load config to get driver
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to load config: %v\n", err)
		os.Exit(1)
	}
	driver := cfg.DB.Connection
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
	// Load config
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to load config: %v\n", err)
		os.Exit(1)
	}
	migrator := cfg.DB.Migrator
	if migrator == "" {
		migrator = "atlas"
	}
	fmt.Println("→ Detected migrator:", migrator)

	if migrator == "gorm" {
		fmt.Println("→ Starting GORM AutoMigrate...")
		// Connect to DB
		db, err := database.Connect(cfg.DB)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Database connection failed: %v\n", err)
			os.Exit(1)
		}

		// Get GORM instance
		gormDB, err := db.GetGorm()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// AutoMigrate
		models := services.Model()
		if err := gormDB.AutoMigrate(models...); err != nil {
			fmt.Fprintf(os.Stderr, "Error: AutoMigrate failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ Database migration completed successfully")
		return
	}

	fmt.Println("→ Running Atlas migrate apply...")
	dbConn := cfg.DB.Connection
	dbName := cfg.DB.Name
	dbUser := cfg.DB.Username
	dbPass := cfg.DB.Password
	dbHost := cfg.DB.Host
	dbPort := fmt.Sprintf("%d", cfg.DB.Port)
	dbDevName := cfg.DB.DevName

	atlasEnv := os.Environ()
	atlasEnv = append(atlasEnv, "DB_CONNECTION="+dbConn)
	atlasEnv = append(atlasEnv, "DB_NAME="+dbName)
	atlasEnv = append(atlasEnv, "DB_USERNAME="+dbUser)
	atlasEnv = append(atlasEnv, "DB_PASSWORD="+dbPass)
	atlasEnv = append(atlasEnv, "DB_HOST="+dbHost)
	atlasEnv = append(atlasEnv, "DB_PORT="+dbPort)
	atlasEnv = append(atlasEnv, "DB_DEV_NAME="+dbDevName)

	cmd := exec.Command("atlas", "migrate", "apply", "--env", "gorm")
	cmd.Env = atlasEnv
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Atlas migrate apply failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Atlas migration completed successfully")
}
