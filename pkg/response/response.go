package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse is a common structure for all JSON responses.
type APIResponse struct {
	Code    int         `json:"code"`           // HTTP status code
	Status  string      `json:"status"`         // "success" or "error"
	Message string      `json:"message"`        // Human-readable status or error message
	Data    interface{} `json:"data,omitempty"` // Omit if null
}

// Success sends a successful JSON response.
func Success(c *gin.Context, httpCode int, message string, data interface{}) {
	// If no specific HTTP code is given, default to 200
	if httpCode == 0 {
		httpCode = http.StatusOK
	}
	// Build the response
	res := APIResponse{
		Code:    httpCode,
		Status:  "success",
		Message: message,
		Data:    data,
	}
	c.JSON(httpCode, res)
}

// Error sends an error JSON response.
func Error(c *gin.Context, httpCode int, message string) {
	// If no specific HTTP code is given, default to 400
	if httpCode == 0 {
		httpCode = http.StatusBadRequest
	}
	// Build the response
	res := APIResponse{
		Code:    httpCode,
		Status:  "error",
		Message: message,
	}
	c.JSON(httpCode, res)
}
