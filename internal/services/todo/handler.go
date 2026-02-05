package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TodoHandler handles HTTP requests for todo resources
type TodoHandler struct {
	// Add repository or service dependencies here
}

// GetAllTodos returns a list of all todos
func (h *TodoHandler) GetAllTodos(c *gin.Context) {
	// Sample data for demonstration
	todos := []Todo{
		{ID: 1, Title: "Learn GoForge", Description: "Start building amazing microservices"},
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Todos retrieved successfully",
		"data":    todos,
	})
}

// GetTodoByID returns a single todo by its ID
func (h *TodoHandler) GetTodoByID(c *gin.Context) {
	id := c.Param("id")
	// Sample logic
	todo := Todo{ID: 1, Title: "Learn GoForge", Description: "Start building amazing microservices"}

	if id != "1" {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Todo not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Todo retrieved successfully",
		"data":    todo,
	})
}

// CreateTodo creates a new todo item
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Logic to save todo
	todo := Todo{
		ID:          2,
		Title:       input.Title,
		Description: input.Description,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Todo created successfully",
		"data":    todo,
	})
}

// UpdateTodo updates an existing todo item
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Completed   bool   `json:"completed"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Logic to update todo by ID
	_ = id

	c.JSON(http.StatusOK, gin.H{
		"message": "Todo updated successfully",
	})
}

// DeleteTodo deletes a todo item by ID
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id := c.Param("id")
	// Logic to delete todo
	_ = id

	c.JSON(http.StatusOK, gin.H{
		"message": "Todo deleted successfully",
	})
}
