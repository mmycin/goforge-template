package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/mmycin/goforge/core/errors"
)

// ErrorHandler is a Gin middleware that catches any error set on the context
// via c.Error() and writes a standardized JSON response.
//
// Usage in handlers:
//
//	if err != nil {
//	    _ = c.Error(apperrors.ErrNotFound("todo not found").Wrap(err))
//	    return
//	}
//
// This middleware must be placed AFTER Recovery() so panics are already
// converted to errors before we inspect them.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Nothing to do if no errors were attached or a response was already sent.
		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}

		// Use the last error (most specific) to build the response.
		lastErr := c.Errors.Last().Err

		var appErr *apperrors.AppError
		if errors.As(lastErr, &appErr) {
			c.JSON(appErr.HTTPStatus, appErr)
			return
		}

		// Fallback for unexpected errors
		c.JSON(http.StatusInternalServerError, &apperrors.AppError{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "an unexpected error occurred",
		})
	}
}
