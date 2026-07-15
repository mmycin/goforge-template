package todo

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	apperrors "github.com/mmycin/goforge/core/errors"
	"github.com/mmycin/goforge/core/validator"
)

type TodoHandler struct {
	service *TodoService
}

func NewTodoHandler(service *TodoService) *TodoHandler {
	return &TodoHandler{service: service}
}

// GetAll returns every todo item.
func (h *TodoHandler) GetAll(c *gin.Context) {
	todos, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		_ = c.Error(apperrors.ErrInternalServer("failed to retrieve todos").Wrap(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data retrieved All",
		"data":    todos,
	})
}

// GetByID returns a single todo by its numeric ID.
func (h *TodoHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		_ = c.Error(apperrors.ErrBadRequest("id must be a positive integer").Wrap(err))
		return
	}

	todo, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		_ = c.Error(apperrors.ErrNotFound("todo not found").Wrap(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Detail retrieved",
		"data":    todo,
	})
}

// CreateTodoRequest is the expected JSON body for creating a todo.
type CreateTodoRequest struct {
	Title       string `json:"title"       validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"required,min=1,max=255"`
}

// Create creates a new todo item.
func (h *TodoHandler) Create(c *gin.Context) {
	var req CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperrors.ErrBadRequest("invalid request body").Wrap(err))
		return
	}

	// Validate using the dedicated validator package
	if appErr := validator.Validate(req); appErr != nil {
		_ = c.Error(appErr)
		return
	}

	todo, err := h.service.Create(c.Request.Context(), req.Title, req.Description)
	if err != nil {
		_ = c.Error(apperrors.ErrInternalServer("failed to create todo").Wrap(err))
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Todo created",
		"data":    todo,
	})
}
