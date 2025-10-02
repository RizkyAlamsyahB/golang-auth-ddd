package response

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Original functions for non-Gin handlers (if needed)
func JSON(w http.ResponseWriter, statusCode int, success bool, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Success: success,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

func Success(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	JSON(w, statusCode, true, message, data)
}

func Error(w http.ResponseWriter, statusCode int, message string, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Success: false,
		Message: message,
		Error:   err,
	}

	json.NewEncoder(w).Encode(response)
}

// New Gin-specific functions
func SuccessGin(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorGin(c *gin.Context, statusCode int, message string, err interface{}) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}
