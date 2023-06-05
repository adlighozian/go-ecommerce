package response

import (
	"net/http"
	"product-colors-go/model"

	"github.com/gin-gonic/gin"
)

func ResponseSuccess(ctx *gin.Context, status int, data any) {
	ctx.JSON(status, model.ResponSuccess{
		Status:  status,
		Message: http.StatusText(status),
		Data:    data,
	})
}

func ResponseError(ctx *gin.Context, status int, err error) {
	ctx.JSON(status, model.ResponError{
		Status:  status,
		Message: http.StatusText(status),
		Error:   err.Error(),
	})
}
