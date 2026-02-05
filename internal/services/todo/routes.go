package todo

import (
	"github.com/gin-gonic/gin"
	"github.com/mmycin/goforge/internal/server"
)

func init() {
	server.Register(&TodoRouter{})
}

// TodoRouter handles routing for the todo service
type TodoRouter struct{}

// Register registers todo routes to the gin engine
func (r *TodoRouter) Register(engine gin.IRouter) {
	h := &TodoHandler{}

	// Pattern: router.Group("/todos", func(route) { ... })
	// We use a closure to match the spirit of the DX request
	RegisterGroup(engine, "/todos", func(group *gin.RouterGroup) {
		group.GET("/", h.GetAllTodos)
		group.GET("/:id", h.GetTodoByID)
		group.POST("/", h.CreateTodo)
		group.PUT("/:id", h.UpdateTodo)
		group.DELETE("/:id", h.DeleteTodo)
	})
}

// RegisterGroup is a helper to support the callback-based grouping DX
func RegisterGroup(engine gin.IRouter, Path string, fn func(*gin.RouterGroup)) {
	group := engine.Group(Path)
	fn(group)
}
