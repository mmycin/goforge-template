package console

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var makeServiceCmd = &cobra.Command{
	Use:   "make:service [name]",
	Short: "Create a new service",
	Long:  `Generate a new service with handler, repository, model, routes, and proto files.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		fmt.Printf("Creating service: %s\n", name)
		makeService(name)
	},
}

func makeService(name string) {
	targetDir := filepath.Join("internal/services", name)

	if _, err := os.Stat(targetDir); err == nil {
		fmt.Printf("Error: Service '%s' already exists\n", name)
		os.Exit(1)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		fmt.Printf("Error: Failed to create directory: %v\n", err)
		os.Exit(1)
	}

	files := map[string]string{
		"handler.go": fmt.Sprintf(`package %s

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type %sHandler struct {}

func (h *%sHandler) GetAll(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Data retrieved",
		"data":    nil,
	})
}

func (h *%sHandler) GetByID(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Detail retrieved",
		"data":    nil,
	})
}
`, name, strings.Title(name), strings.Title(name), strings.Title(name)),
		"proto.go": "package " + name + "\n",
		"repo.go":  "package " + name + "\n",
		"routes.go": fmt.Sprintf(`package %s

import (
	"github.com/gin-gonic/gin"
	"github.com/mmycin/goforge/internal/server"
)

func init() {
	server.Register(&%sRouter{})
}

type %sRouter struct{}

func (r *%sRouter) Register(engine gin.IRouter) {
	h := &%sHandler{}
	group := engine.Group("/%ss")
	{
		group.GET("", h.GetAll)
		group.GET("/:id", h.GetByID)
	}
}
`, name, strings.Title(name), strings.Title(name), strings.Title(name), strings.Title(name), name),
		"service.go": "package " + name + "\n",
		"model.go":   fmt.Sprintf("package %s\n\nimport \"time\"\n\ntype %s struct {\n\tID        uint      `gorm:\"primaryKey;autoIncrement\"`\n\tCreatedAt time.Time `gorm:\"autoCreateTime\"`\n\tUpdatedAt time.Time `gorm:\"autoUpdateTime\"`\n}\n", name, strings.Title(name)),
	}

	for fname, content := range files {
		if err := os.WriteFile(filepath.Join(targetDir, fname), []byte(content), 0644); err != nil {
			fmt.Printf("Error: Failed to write %s: %v\n", fname, err)
		}
	}

	updateKernel(name)

	fmt.Printf("✓ Service '%s' created successfully and auto-registered\n", name)
}

func updateKernel(name string) {
	kernelPath := "internal/services/kernel.go"
	importLine := fmt.Sprintf("\t_ \"github.com/mmycin/goforge/internal/services/%s\"", name)

	content, err := os.ReadFile(kernelPath)
	if err != nil {
		fmt.Printf("Warning: Could not read %s: %v\n", kernelPath, err)
		return
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	added := false

	for _, line := range lines {
		if !added && strings.Contains(line, ")") && strings.HasPrefix(strings.TrimSpace(line), ")") {
			newLines = append(newLines, importLine)
			added = true
		}
		newLines = append(newLines, line)
	}

	if !added {
		fmt.Printf("Warning: Could not automatically update %s. Please add %s manually.\n", kernelPath, importLine)
		return
	}

	err = os.WriteFile(kernelPath, []byte(strings.Join(newLines, "\n")), 0644)
	if err != nil {
		fmt.Printf("Warning: Could not write to %s: %v\n", kernelPath, err)
	}
}
