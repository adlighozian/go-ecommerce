package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResponseSuccess(ctx *gin.Context, status int, data any) {
	ctx.JSON(status, map[string]interface{}{
		"data":    data,
		"message": http.StatusText(status),
		"status":  status,
	})
}

func ResponseError(ctx *gin.Context, status int, err error) {
	ctx.JSON(status, map[string]interface{}{
		"error":   err.Error(),
		"message": http.StatusText(status),
		"status":  status,
	})
}
