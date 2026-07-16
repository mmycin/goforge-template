package errors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Respond writes the error as a JSON response to the Gin context.
// If err is an *AppError it uses its HTTP status and body.
// All other errors fall back to a generic 500 response so internal
// details are never leaked to the client.
func Respond(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.HTTPStatus, appErr)
		return
	}

	c.JSON(http.StatusInternalServerError, &AppError{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "an unexpected error occurred",
	})
}
