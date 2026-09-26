package console

import (
	"fmt"
	"os"
	"os/exec"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/mmycin/goforge/core/config"
	"github.com/mmycin/goforge/core/database"
	"github.com/mmycin/goforge/internal/services"
	"github.com/spf13/cobra"
)

// atlasConfig is the Atlas --config URI.
// Atlas requires the file:// scheme regardless of OS.
const atlasConfig = "file://core/database/atlas.hcl"

// sqlcConfig is the path to the sqlc config file relative to the project root.
const sqlcConfig = "core/database/sqlc.yaml"

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

		db, err := database.Connect(cfg.DB)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Database connection failed: %v\n", err)
			os.Exit(1)
		}

		gormDB, err := db.GetGorm()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		models := services.Model()
		if err := gormDB.AutoMigrate(models...); err != nil {
			fmt.Fprintf(os.Stderr, "Error: AutoMigrate failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ Database migration completed successfully")
		return
	}

	fmt.Println("→ Running Atlas migrate apply...")

	atlasEnv := os.Environ()
	atlasEnv = append(atlasEnv, "DB_CONNECTION="+cfg.DB.Connection)
	atlasEnv = append(atlasEnv, "DB_NAME="+cfg.DB.Name)
	atlasEnv = append(atlasEnv, "DB_USERNAME="+cfg.DB.Username)
	atlasEnv = append(atlasEnv, "DB_PASSWORD="+cfg.DB.Password)
	atlasEnv = append(atlasEnv, "DB_HOST="+cfg.DB.Host)
	atlasEnv = append(atlasEnv, "DB_PORT="+fmt.Sprintf("%d", cfg.DB.Port))
	atlasEnv = append(atlasEnv, "DB_DEV_NAME="+cfg.DB.DevName)

	// --config points Atlas at core/database/atlas.hcl.
	// All paths inside that file are still relative to CWD (the project root).
	cmd := exec.Command("atlas", "migrate", "apply",
		"--env", "gorm",
		"--config", atlasConfig,
	)
	cmd.Env = atlasEnv
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Atlas migrate apply failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Atlas migration completed successfully")
}
