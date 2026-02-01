package console

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mmycin/goforge/internal/database"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func makeMigration(name string) {
	// 1. Register Models
	log.Println("Registering models...")
	if err := registerModels(); err != nil {
		log.Fatalf("Failed to register models: %v", err)
	}

	// 2. Atlas Migrate Hash (before Diff)
	log.Println("Running atlas migrate hash...")
	cmd := exec.Command("atlas", "migrate", "hash", "--env", "gorm")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Atlas hash failed: %v", err)
	}

	// 3. Atlas Migrate Diff
	log.Println("Running atlas migrate diff...")
	cmd = exec.Command("atlas", "migrate", "diff", "--env", "gorm", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Atlas migration failed: %v", err)
	}

	// 4. Cleanup SQL (sed replacement)
	log.Println("Applying sed fix...")
	files, _ := filepath.Glob("internal/database/migrations/*.sql")
	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			log.Printf("Failed to read %s: %v", f, err)
			continue
		}
		newContent := strings.ReplaceAll(string(content), "`", "")
		if err := os.WriteFile(f, []byte(newContent), 0644); err != nil {
			log.Printf("Failed to write %s: %v", f, err)
		}
	}

	// 5. SQLC Generate
	makeGen() // Reuse the gen command logic

	log.Println("Migration creation flow completed.")
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

func migrateDB() {
	dbName := os.Getenv("DB_NAME")
	dbDriver := os.Getenv("DB_DRIVER") // mysql, postgres, sqlite, sqlserver
	dbDsn := os.Getenv("DB_DSN")

	if dbDriver == "" {
		// Fallback or guess
		if strings.HasSuffix(dbName, ".db") {
			dbDriver = "sqlite"
			if dbDsn == "" {
				dbDsn = dbName
			}
		} else {
			// default to sqlite if nothing specified, or check DB_CONNECTION
			conn := os.Getenv("DB_CONNECTION") // Laravel style
			if conn != "" {
				dbDriver = conn
			} else {
				dbDriver = "sqlite" // default
			}
		}
	}

	if dbDsn == "" {
		dbDsn = dbName // hope for the best
	}

	log.Printf("Connecting to DB: %s (%s)", dbDriver, dbDsn)

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
		log.Fatalf("Unsupported driver: %s", dbDriver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Running GORM AutoMigrate...")
	models := database.Model()
	if err := db.AutoMigrate(models...); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	log.Println("Migration completed.")
}
