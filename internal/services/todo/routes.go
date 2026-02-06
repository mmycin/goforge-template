package todo

import (
	"github.com/gin-gonic/gin"
	"github.com/mmycin/goforge/internal/server"
)

func init() {
	server.Register(&TodoRoutes{})
}

type TodoRoutes struct{}

func (r *TodoRoutes) Register(engine *gin.Engine) {
	h := &TodoHandler{}

	group := engine.Group("/api/todos")
	// group.Use(server.AuthMiddleware())

	group.GET("/", h.GetAll)
	group.GET("/:id", h.GetByID)
}
