package todo

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type TodoHandler struct{}

func (h *TodoHandler) GetAll(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Data retrieved All",
		"data":    "Dummy Todo",
	})
}

func (h *TodoHandler) GetByID(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Detail retrieved",
		"data":    "Dummy Todo",
	})
}
