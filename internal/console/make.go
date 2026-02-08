package console

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mmycin/goforge/internal/config"
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
	moduleName := config.App.Module
	if moduleName == "" {
		// Fallback or error if not set, though Load() should handle it
		moduleName = "github.com/mmycin/goforge"
		Warning("APP_MODULE not set, using default: %s", moduleName)
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
		"data":    []string{},
	})
}

func (h *%sHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Detail retrieved",
		"data":    id,
	})
}
`, name, camelName, camelName, camelName),
		"grpc.go": "package " + name + "\n",
		"routes.go": fmt.Sprintf(`package %s

import (
	"github.com/gin-gonic/gin"
	"%s/internal/server"
)

func init() {
	server.Register(&%sRoutes{})
}

type %sRoutes struct{}

func (r *%sRoutes) Register(engine *gin.Engine) {
	h := &%sHandler{}

	group := engine.Group("/api/%ss")
	// Middleware is applied globally in server/http.go

	group.GET("/", h.GetAll)
	group.GET("/:id", h.GetByID)
}
`, name, moduleName, camelName, camelName, camelName, camelName, name),
		"docs.go": fmt.Sprintf(`package %s

import (
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"%s/internal/server"
)

func init() {
	server.Register(&%sDocs{})
}

type %sDocs struct{}

func (d *%sDocs) Register(engine *gin.Engine) {
	config := server.NewHumaConfig("%s API", "1.0.0", "/api/docs/%s")
	
	// Create API instance
	humagin.New(engine, config)
	
	// Register Huma operations here if you want to generate OpenAPI spec
	// usage: huma.Register(api, operation, handler)
}
`, name, moduleName, camelName, camelName, camelName, name, name),
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
		// Proto content stays similar but ensure package name is simple
		protoContent := fmt.Sprintf(`syntax = "proto3";

package %s;

option go_package = "%s/proto/%s/gen";

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
`, name, moduleName, name, camelName, camelName)

		if err := os.WriteFile(filepath.Join(protoDir, name+".proto"), []byte(protoContent), 0644); err != nil {
			Error("Failed to write proto file: %v", err)
		}
	}

	if err := registerModels(); err != nil {
		Warning("Could not automatically update kernel.go: %v", err)
	}

	Success("Service '%s' created successfully and auto-registered in kernel.go", name)
}
