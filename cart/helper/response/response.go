package response

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func ResponseSuccess(ctx *gin.Context, status int, message string, data interface{}) {
	if message == "" {
		message = http.StatusText(status)
	}
	ctx.JSON(status, map[string]interface{}{
		"status":  status,
		"message": message,
		"data"  :  data,
	})
}

func ResponseError(ctx *gin.Context, status int, message string, err error) {
	if message == "" {
		message = http.StatusText(status)
	}
	ctx.JSON(status, map[string]interface{}{
		"status":  status,
		"message": message,
		"error":   err.Error(),
	})
}