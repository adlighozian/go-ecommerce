package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JSONRes struct {
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type JSONGatewayRes struct {
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Data    string `json:"data,omitempty"`
}

func NewJSONRes(c *gin.Context, statusCode int, message string, data any) {
	if message == "" {
		message = http.StatusText(statusCode)
	}
	c.JSON(statusCode, JSONRes{
		Status:  statusCode,
		Message: message,
		Data:    data,
	})
}

func NewJSONResErr(c *gin.Context, statusCode int, message string, err string) {
	if message == "" {
		message = http.StatusText(statusCode)
	}
	c.JSON(statusCode, JSONRes{
		Status:  statusCode,
		Message: message,
		Error:   err,
	})
}
