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
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func makeMigration(name string) {
	// 1. Register Models
	fmt.Println("→ Registering models...")
	if err := registerModels(); err != nil {
		fmt.Printf("Error: Failed to register models: %v\n", err)
		os.Exit(1)
	}

	// 2. Atlas Migrate Hash (before Diff)
	fmt.Println("→ Running atlas migrate hash...")
	cmd := exec.Command("atlas", "migrate", "hash", "--env", "gorm")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: Atlas hash failed: %v\n", err)
		os.Exit(1)
	}

	// 3. Atlas Migrate Diff
	fmt.Println("→ Running atlas migrate diff...")
	cmd = exec.Command("atlas", "migrate", "diff", "--env", "gorm", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: Atlas migration failed: %v\n", err)
		os.Exit(1)
	}

	// 4. Cleanup SQL (sed replacement)
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
	loader := gormschema.New("sqlite")
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
		// Build DSN based on driver
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
