package console

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func makeService(name string) {
	log.Printf("Creating service: %s", name)
	targetDir := filepath.Join("internal/services", name)

	// Check if service already exists
	if _, err := os.Stat(targetDir); err == nil {
		log.Fatalf("Service %s already exists", name)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	// Templates based on 'todo' service.
	files := map[string]string{
		"handler.go": "package " + name + "\n",
		"proto.go":   "package " + name + "\n",
		"repo.go":    "package " + name + "\n",
		"service.go": "package " + name + "\n",
		"model.go":   fmt.Sprintf("package %s\n\nimport \"time\"\n\ntype %s struct {\n\tID        uint      `gorm:\"primaryKey;autoIncrement\"`\n\tCreatedAt time.Time `gorm:\"autoCreateTime\"`\n\tUpdatedAt time.Time `gorm:\"autoUpdateTime\"`\n}\n", name, strings.ToUpper(name[:1])+name[1:]),
	}

	for fname, content := range files {
		if err := os.WriteFile(filepath.Join(targetDir, fname), []byte(content), 0644); err != nil {
			log.Printf("Failed to write %s: %v", fname, err)
		}
	}
	log.Println("Service created.")
}
