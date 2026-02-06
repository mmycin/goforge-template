package console

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/mmycin/goforge/internal/config"
	"github.com/mmycin/goforge/internal/database"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Execute GORM AutoMigrate to synchronize database schema with models.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running database migration...")
		migrateDB()
	},
}

// genMigrationCmd represents the gen:migration command
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

// loaderCmd represents the loader command
var loaderCmd = &cobra.Command{
	Use:   "loader",
	Short: "Run GORM schema loader",
	Long:  `Load and display GORM schema definitions.`,
	Run: func(cmd *cobra.Command, args []string) {
		runLoader()
	},
}

// genSqlcCmd represents the gen:sqlc command
var genSqlcCmd = &cobra.Command{
	Use:   "gen:sqlc",
	Short: "Run code generation",
	Long:  `Execute sqlc generate to create database query code.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running code generation...")
		makeGen()
	},
}

func makeGen() {
	fmt.Println("Executing sqlc generate...")
	cmd := exec.Command("sqlc", "generate")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: sqlc generate failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Code generation completed successfully")
}

func makeMigration(name string) {
	fmt.Println("→ Registering models...")
	if err := registerModels(); err != nil {
		fmt.Printf("Error: Failed to register models: %v\n", err)
		os.Exit(1)
	}

	// Prepare environment for Atlas by propagating database config
	atlasEnv := os.Environ()
	atlasEnv = append(atlasEnv, "DB_CONNECTION="+config.DB.Connection)
	atlasEnv = append(atlasEnv, "DB_NAME="+config.DB.Name)
	atlasEnv = append(atlasEnv, "DB_USERNAME="+config.DB.Username)
	atlasEnv = append(atlasEnv, "DB_PASSWORD="+config.DB.Password)
	atlasEnv = append(atlasEnv, "DB_HOST="+config.DB.Host)
	atlasEnv = append(atlasEnv, "DB_PORT="+fmt.Sprintf("%d", config.DB.Port))
	atlasEnv = append(atlasEnv, "DB_DEV_NAME="+config.DB.DevName)

	fmt.Println("→ Running atlas migrate hash...")
	cmd := exec.Command("atlas", "migrate", "hash", "--env", "gorm")
	cmd.Env = atlasEnv
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: Atlas hash failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("→ Running atlas migrate diff...")
	cmd = exec.Command("atlas", "migrate", "diff", "--env", "gorm", name)
	cmd.Env = atlasEnv
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("\nError: Atlas migration failed: %v\n", err)
		if config.DB.DevName == "" && (config.DB.Connection == "mysql" || config.DB.Connection == "postgres") {
			fmt.Println("\nTIP: Atlas requires a clean/empty database for the 'dev' environment.")
			fmt.Printf("1. Create an empty database in your %s server (e.g., 'CREATE DATABASE %s_dev;')\n", config.DB.Connection, config.DB.Name)
			fmt.Printf("2. Add 'DB_DEV_NAME=%s_dev' to your .env file\n", config.DB.Name)
			fmt.Println("3. Run the command again.")
		}
		os.Exit(1)
	}

	fmt.Println("→ Cleaning up SQL files...")
	files, _ := filepath.Glob("internal/database/migrations/*.sql")
	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			fmt.Printf("Warning: Failed to read %s: %v\n", f, err)
			continue
		}
		newContent := strings.ReplaceAll(string(content), "`", "")
		if err := os.WriteFile(f, []byte(newContent), 0644); err != nil {
			fmt.Printf("Warning: Failed to write %s: %v\n", f, err)
		}
	}
	fmt.Println("✓ Migration created successfully")
}

func registerModels() error {
	servicesDir := "internal/services"
	entries, err := os.ReadDir(servicesDir)
	if err != nil {
		return err
	}

	var services []string
	for _, e := range entries {
		if e.IsDir() {
			modelPath := filepath.Join(servicesDir, e.Name(), "model.go")
			if _, err := os.Stat(modelPath); err == nil {
				services = append(services, e.Name())
			}
		}
	}

	tmpl := `package database

import (
{{- range . }}
	"github.com/mmycin/goforge/internal/services/{{ . }}"
{{- end }}
)

func Model() []any {
	return []any{
{{- range . }}
		&{{ . }}.{{ title . }}{},
{{- end }}
	}
}
`
	funcMap := template.FuncMap{
		"title": func(s string) string {
			if len(s) == 0 {
				return ""
			}
			return strings.ToUpper(s[:1]) + s[1:]
		},
	}

	t, err := template.New("kernel").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return err
	}

	f, err := os.Create("internal/database/kernel.go")
	if err != nil {
		return err
	}
	defer f.Close()

	return t.Execute(f, services)
}

func runLoader() {
	models := database.Model()
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

func migrateDB() {
	dbName := config.DB.Name
	dbDriver := config.DB.Connection
	dbDsn := ""

	if dbDriver == "" {
		if strings.HasSuffix(dbName, ".db") {
			dbDriver = "sqlite"
			dbDsn = dbName
		} else {
			dbDriver = "sqlite"
		}
	}

	if dbDsn == "" {
		switch dbDriver {
		case "mysql":
			dbDsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				config.DB.Username, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.Name)
		case "postgres":
			dbDsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
				config.DB.Host, config.DB.Username, config.DB.Password, config.DB.Name, config.DB.Port)
		case "sqlite":
			dbDsn = dbName
		case "sqlserver":
			dbDsn = fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
				config.DB.Username, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.Name)
		}
	}

	fmt.Printf("→ Connecting to database: %s\n", dbDriver)

	var dialector gorm.Dialector
	switch dbDriver {
	case "mysql":
		dialector = mysql.Open(dbDsn)
	case "postgres":
		dialector = postgres.Open(dbDsn)
	case "sqlite":
		dialector = sqlite.Open(dbDsn)
	case "sqlserver":
		dialector = sqlserver.Open(dbDsn)
	default:
		fmt.Printf("Error: Unsupported database driver: %s\n", dbDriver)
		os.Exit(1)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		fmt.Printf("Error: Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("→ Running GORM AutoMigrate...")
	models := database.Model()
	if err := db.AutoMigrate(models...); err != nil {
		fmt.Printf("Error: AutoMigrate failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Database migration completed successfully")
}
