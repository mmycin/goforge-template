package console

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func makeService(name string) {
	targetDir := filepath.Join("internal/services", name)

	// Check if service already exists
	if _, err := os.Stat(targetDir); err == nil {
		fmt.Printf("Error: Service '%s' already exists\n", name)
		os.Exit(1)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		fmt.Printf("Error: Failed to create directory: %v\n", err)
		os.Exit(1)
	}

	// Templates based on 'todo' service.
	files := map[string]string{
		"handler.go": "package " + name + "\n",
		"proto.go":   "package " + name + "\n",
		"repo.go":    "package " + name + "\n",
		"routes.go":  "package " + name + "\n",
		"service.go": "package " + name + "\n",
		"model.go":   fmt.Sprintf("package %s\n\nimport \"time\"\n\ntype %s struct {\n\tID        uint      `gorm:\"primaryKey;autoIncrement\"`\n\tCreatedAt time.Time `gorm:\"autoCreateTime\"`\n\tUpdatedAt time.Time `gorm:\"autoUpdateTime\"`\n}\n", name, strings.ToUpper(name[:1])+name[1:]),
	}

	for fname, content := range files {
		if err := os.WriteFile(filepath.Join(targetDir, fname), []byte(content), 0644); err != nil {
			fmt.Printf("Error: Failed to write %s: %v\n", fname, err)
		}
	}
	fmt.Printf("✓ Service '%s' created successfully\n", name)
}
