package todo

import (
	"github.com/gin-gonic/gin"
)

type TodoRoutes struct {
	handler *TodoHandler
}

func NewTodoRoutes(handler *TodoHandler) *TodoRoutes {
	return &TodoRoutes{handler: handler}
}

func (r *TodoRoutes) Register(engine *gin.Engine) {
	group := engine.Group("/api/todos")
	// group.Use(server.AuthMiddleware())

	group.GET("/", r.handler.GetAll)
	group.GET("/:id", r.handler.GetByID)
	group.POST("/", r.handler.Create)
}
