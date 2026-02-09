package todo

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
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
	api := humagin.New(engine, config)

	// Register health check
	huma.Register(api, huma.Operation{
		OperationID: "get-health",
		Method:      http.MethodGet,
		Path:        "/api/todo/health",
		Summary:     "Health check",
		Description: "Check if the service is healthy",
		Tags:        []string{"Health"},
	}, func(ctx context.Context, input *struct{}) (*struct{ Body string }, error) {
		return &struct{ Body string }{Body: "OK"}, nil
	})
}
