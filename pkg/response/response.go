package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success sends a success JSON response.
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, APIResponse{
		Message: "success",
		Data:    data,
	})
}

// Error sends an error JSON response.
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, APIResponse{
		Message: message,
		Data:    nil,
	})
}
