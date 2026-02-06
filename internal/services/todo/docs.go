package todo

import (
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/mmycin/goforge/internal/server"
)

func init() {
	server.Register(&TodoDocs{})
}

type TodoDocs struct{}

func (d *TodoDocs) Register(engine *gin.Engine) {
	// Use the global config helper to get pre-configured Huma config with security
	config := server.NewHumaConfig("Todo API", "1.0.0", "/api/docs/todo")

	// Create API instance to serve docs
	humagin.New(engine, config)

	// Register Huma operations here if you want to generate OpenAPI spec
}
