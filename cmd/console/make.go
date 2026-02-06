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
		Info("Creating service: %s", name)
		makeService(name)
	},
}

func makeService(name string) {
	targetDir := filepath.Join("internal/services", name)

	if _, err := os.Stat(targetDir); err == nil {
		Error("Service '%s' already exists", name)
		os.Exit(1)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		Error("Failed to create directory: %v", err)
		os.Exit(1)
	}

	camelName := toCamelCase(name)

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
`, name, camelName, camelName, camelName),
		"grpc.go": "package " + name + "\n",
		"repo.go": "package " + name + "\n",
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
`, name, camelName, camelName, camelName, camelName, name),
		"service.go": "package " + name + "\n",
		"model.go":   fmt.Sprintf("package %s\n\nimport \"time\"\n\ntype %s struct {\n\tID        uint      `gorm:\"primaryKey;autoIncrement\"`\n\tCreatedAt time.Time `gorm:\"autoCreateTime\"`\n\tUpdatedAt time.Time `gorm:\"autoUpdateTime\"`\n}\n", name, camelName),
	}

	for fname, content := range files {
		if err := os.WriteFile(filepath.Join(targetDir, fname), []byte(content), 0644); err != nil {
			Error("Failed to write %s: %v", fname, err)
		}
	}

	// Create proto directory and file
	protoDir := filepath.Join("proto", name)
	if err := os.MkdirAll(protoDir, 0755); err != nil {
		Error("Failed to create proto directory: %v", err)
	} else {
		protoContent := fmt.Sprintf(`syntax = "proto3";

package %s;

option go_package = "github.com/mmycin/goforge/proto/%s/gen";

service %sService {
	rpc Create(CreateRequest) returns (CreateResponse);
	rpc Get(GetRequest) returns (GetResponse);
	rpc List(ListRequest) returns (ListResponse);
	rpc Update(UpdateRequest) returns (UpdateResponse);
	rpc Delete(DeleteRequest) returns (DeleteResponse);
}

message %s {
	string id = 1;
	string created_at = 2;
	string updated_at = 3;
}

message CreateRequest {}
message CreateResponse {}

message GetRequest { string id = 1; }
message GetResponse {}

message ListRequest { int32 page = 1; int32 limit = 2; }
message ListResponse {}

message UpdateRequest { string id = 1; }
message UpdateResponse {}

message DeleteRequest { string id = 1; }
message DeleteResponse {}
`, name, name, camelName, camelName)

		if err := os.WriteFile(filepath.Join(protoDir, name+".proto"), []byte(protoContent), 0644); err != nil {
			Error("Failed to write proto file: %v", err)
		}
	}

	updateKernel(name)

	Success("Service '%s' created successfully with proto template and auto-registered", name)
}

func updateKernel(name string) {
	kernelPath := "internal/services/kernel.go"
	importLine := fmt.Sprintf("\t_ \"github.com/mmycin/goforge/internal/services/%s\"", name)

	content, err := os.ReadFile(kernelPath)
	if err != nil {
		Warning("Could not read %s: %v", kernelPath, err)
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
		Warning("Could not automatically update %s. Please add %s manually.", kernelPath, importLine)
		return
	}

	err = os.WriteFile(kernelPath, []byte(strings.Join(newLines, "\n")), 0644)
	if err != nil {
		Warning("Could not write to %s: %v", kernelPath, err)
	}
}
