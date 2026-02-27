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
	"github.com/mmycin/goforge/internal/services"
	"github.com/spf13/cobra"
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

// remSqlcCmd represents the rem:sqlc command
var remSqlcCmd = &cobra.Command{
	Use:   "rem:sqlc",
	Short: "Remove SQLC integration",
	Long:  `Remove generated SQLC code and revert database kernel integration.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Removing SQLC integration...")
		removeSqlc("internal/database/database.go")
	},
}

func makeGen() {
	if err := updateSqlcConfig(); err != nil {
		fmt.Printf("Warning: Failed to update sqlc.yaml: %v\n", err)
	}

	engine := config.DB.Connection
	if engine == "postgres" || engine == "postgresql" {
		engine = "postgresql"
	} else if engine == "mysql" {
		engine = "mysql"
	} else {
		engine = "sqlite"
	}

	fmt.Printf("→ Transforming queries for %s engine...\n", engine)
	if err := transformQueries(engine); err != nil {
		fmt.Printf("Warning: Failed to transform queries: %v\n", err)
	}

	fmt.Println("Executing sqlc generate...")
	cmd := exec.Command("sqlc", "generate")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: sqlc generate failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Code generation completed successfully")

	injectSqlc("internal/database/database.go")
}

func transformQueries(engine string) error {
	queriesDir := "internal/database/queries"
	files, err := filepath.Glob(filepath.Join(queriesDir, "*.sql"))
	if err != nil {
		return err
	}

	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			return err
		}

		newContent := string(content)
		if engine == "postgresql" {
			// Convert ? to $1, $2, etc. (per query basis)
			lines := strings.Split(newContent, "\n")
			placeholderIdx := 1
			for i, line := range lines {
				if strings.Contains(line, "-- name:") {
					placeholderIdx = 1
				}
				for strings.Contains(lines[i], "?") {
					lines[i] = strings.Replace(lines[i], "?", fmt.Sprintf("$%d", placeholderIdx), 1)
					placeholderIdx++
				}
			}
			newContent = strings.Join(lines, "\n")
		} else if engine == "mysql" || engine == "sqlite" {
			// Convert $n to ?
			for i := 1; i < 50; i++ {
				newContent = strings.ReplaceAll(newContent, fmt.Sprintf("$%d", i), "?")
			}
		}

		if err := os.WriteFile(f, []byte(newContent), 0644); err != nil {
			return err
		}
	}
	return nil
}

func updateSqlcConfig() error {
	configPath := "sqlc.yaml"
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	engine := config.DB.Connection
	if engine == "postgres" {
		engine = "postgresql"
	}
	if engine == "" {
		engine = "sqlite"
	}

	lines := strings.Split(string(content), "\n")
	updated := false
	for i, line := range lines {
		if strings.Contains(line, "engine:") {
			// Preserving indentation and comments
			parts := strings.SplitN(line, "engine:", 2)
			if len(parts) == 2 {
				indent := parts[0]
				suffix := ""
				if idx := strings.Index(parts[1], "#"); idx != -1 {
					suffix = " " + parts[1][idx:]
				}
				lines[i] = fmt.Sprintf("%sengine: %q%s", indent, engine, suffix)
				updated = true
				break
			}
		}
	}

	if !updated {
		return fmt.Errorf("could not find 'engine' field in %s", configPath)
	}

	return os.WriteFile(configPath, []byte(strings.Join(lines, "\n")), 0644)
}

func injectSqlc(targetPath string) {
	content, err := os.ReadFile(targetPath)
	if err != nil {
		fmt.Printf("Warning: Could not read %s for injection: %v\n", targetPath, err)
		return
	}

	code := string(content)
	if strings.Contains(code, "sqlc.New(sqlDB)") {
		return // Already injected
	}

	fmt.Println("→ Injecting SQLC support into database...")

	// 1. Inject Import
	if !strings.Contains(code, "internal/database/gen") {
		code = strings.Replace(code,
			fmt.Sprintf("\t\"%s/internal/config\"", config.App.Module),
			fmt.Sprintf("\t\"%s/internal/config\"\n\tsqlc \"%s/internal/database/gen\"", config.App.Module, config.App.Module), 1)
	}

	// 2. Inject Field
	if !strings.Contains(code, "Sqlc *sqlc.Queries") {
		code = strings.Replace(code,
			"// Sqlc field will be added when generated code is available",
			"Sqlc *sqlc.Queries", 1)
	}

	// 3. Inject Initialization
	if !strings.Contains(code, "sqlc.New(sqlDB)") {
		oldInit := "DB = &Database{\n\t\tGorm: gormDB,\n\t}"
		// Try to find it more flexibly if the above exact match fails
		if !strings.Contains(code, oldInit) {
			// fallback to just the struct assignment if exact whitespace fails
			code = strings.Replace(code, "Gorm: gormDB,", "Gorm: gormDB,\n\t\tSqlc: sqlc.New(sqlDB),", 1)
			if !strings.Contains(code, "sqlDB, err := gormDB.DB()") {
				code = strings.Replace(code, "gormDB, err := gorm.Open(dialector, &gorm.Config{})",
					"gormDB, err := gorm.Open(dialector, &gorm.Config{})\n\tif err != nil {\n\t\treturn err\n\t}\n\n\tsqlDB, err := gormDB.DB()", 1)
			}
		} else {
			newInit := `	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	DB = &Database{
		Gorm: gormDB,
		Sqlc: sqlc.New(sqlDB),
	}`
			code = strings.Replace(code, oldInit, newInit, 1)
		}
	}

	if err := os.WriteFile(targetPath, []byte(code), 0644); err != nil {
		fmt.Printf("Warning: Failed to inject SQLC support: %v\n", err)
	}
	fmt.Println("✓ SQLC support injected into database")
}

func removeSqlc(targetPath string) {
	content, err := os.ReadFile(targetPath)
	if err != nil {
		fmt.Printf("Warning: Could not read %s for removal: %v\n", targetPath, err)
		return
	}

	code := string(content)

	// Check if there's any SQLC integration to remove
	hasSqlcImport := strings.Contains(code, "internal/database/gen")
	hasSqlcField := strings.Contains(code, "Sqlc *sqlc.Queries")
	hasSqlcInit := strings.Contains(code, "sqlc.New(sqlDB)")

	if !hasSqlcImport && !hasSqlcField && !hasSqlcInit {
		fmt.Println("No SQLC integration found to remove.")
		return
	}

	fmt.Println("→ Removing SQLC support from database...")

	// 1. Remove Import (handle multiple possible formats)
	if hasSqlcImport {
		// Try with newline and tab prefix
		code = strings.Replace(code, fmt.Sprintf("\n\tsqlc \"%s/internal/database/gen\"", config.App.Module), "", 1)
		// Try with just tab prefix (in case it's the last import)
		code = strings.Replace(code, fmt.Sprintf("\tsqlc \"%s/internal/database/gen\"\n", config.App.Module), "", 1)
		// Try standalone line
		code = strings.Replace(code, fmt.Sprintf("sqlc \"%s/internal/database/gen\"\n", config.App.Module), "", 1)
	}

	// 2. Revert Field (if it exists)
	if hasSqlcField {
		code = strings.Replace(code,
			"Sqlc *sqlc.Queries",
			"// Sqlc field will be added when generated code is available", 1)
	}

	// 3. Revert Initialization (if it exists)
	if hasSqlcInit {
		oldInit := `	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	DB = &Database{
		Gorm: gormDB,
		Sqlc: sqlc.New(sqlDB),
	}`
		newInit := `	DB = &Database{
		Gorm: gormDB,
	}`
		code = strings.Replace(code, oldInit, newInit, 1)
	}

	if err := os.WriteFile(targetPath, []byte(code), 0644); err != nil {
		fmt.Printf("Warning: Failed to remove SQLC support: %v\n", err)
		return
	}

	// Optional: Remove generated files?
	// os.RemoveAll("internal/database/gen")

	fmt.Println("✓ SQLC support removed from database")
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

	fmt.Println("→ Running atlas migrate diff...")
	cmd := exec.Command("atlas", "migrate", "diff", "--env", "gorm", name)
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
		newContent := string(content)
		newContent = strings.ReplaceAll(newContent, "`", "")

		// Clean up MySQL specific items that break SQLC (even in MySQL mode sometimes)
		// Remove COLLATE ...
		lines := strings.Split(newContent, "\n")
		for i, line := range lines {
			if idx := strings.Index(line, "COLLATE"); idx != -1 {
				lines[i] = strings.TrimSpace(line[:idx])
				if strings.HasSuffix(lines[i], ";") {
					// Keep the semicolon
				} else {
					lines[i] += ";"
				}
			}
		}
		newContent = strings.Join(lines, "\n")

		if err := os.WriteFile(f, []byte(newContent), 0644); err != nil {
			fmt.Printf("Warning: Failed to write %s: %v\n", f, err)
		}
	}

	fmt.Println("→ Running atlas migrate hash...")
	cmd = exec.Command("atlas", "migrate", "hash", "--env", "gorm")
	cmd.Env = atlasEnv
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: Atlas hash failed: %v\n", err)
		os.Exit(1)
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

	tmpl := `package services

import (
	"{{ .Module }}/internal/server"
{{- range .Services }}
	"{{ $.Module }}/internal/services/{{ . }}"
{{- end }}
)

// GetRouters returns all service routers to be registered
func GetRouters() []server.Router {
	return server.GetRegisteredRouters()
}

// Model returns all models to be registered with GORM
func Model() []any {
	return []any{
{{- range .Services }}
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

	f, err := os.Create("internal/services/kernel.go")
	if err != nil {
		return err
	}
	defer f.Close()

	data := struct {
		Module   string
		Services []string
	}{
		Module:   config.App.Module,
		Services: services,
	}

	return t.Execute(f, data)
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

func migrateDB() {
	fmt.Println("→ Connecting to database...")
	if err := database.Connect(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	migrator := strings.ToLower(config.DB.Migrator)
	if migrator == "atlas" {
		fmt.Println("→ Running Atlas migrate apply...")
		runAtlasMigrate()
	} else {
		fmt.Println("→ Running GORM AutoMigrate...")
		models := services.Model()
		if err := database.DB.Gorm.AutoMigrate(models...); err != nil {
			fmt.Printf("Error: AutoMigrate failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ Database migration completed successfully")
	}
}

func runAtlasMigrate() {
	// Prepare environment for Atlas by propagating database config
	atlasEnv := os.Environ()
	atlasEnv = append(atlasEnv, "DB_CONNECTION="+config.DB.Connection)
	atlasEnv = append(atlasEnv, "DB_NAME="+config.DB.Name)
	atlasEnv = append(atlasEnv, "DB_USERNAME="+config.DB.Username)
	atlasEnv = append(atlasEnv, "DB_PASSWORD="+config.DB.Password)
	atlasEnv = append(atlasEnv, "DB_HOST="+config.DB.Host)
	atlasEnv = append(atlasEnv, "DB_PORT="+fmt.Sprintf("%d", config.DB.Port))
	atlasEnv = append(atlasEnv, "DB_DEV_NAME="+config.DB.DevName)

	cmd := exec.Command("atlas", "migrate", "apply", "--env", "gorm")
	cmd.Env = atlasEnv
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: Atlas migrate apply failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Atlas migration completed successfully")
}
