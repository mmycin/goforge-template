package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct{}

func (h *TodoHandler) GetAll(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Data retrieved",
		"data":    []string{},
	})
}

func (h *TodoHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Detail retrieved",
		"data":    id,
	})
}
