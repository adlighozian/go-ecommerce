package response

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Status	int			`json:"status,omniempty"`
	Message string		`json:"message,omniempty"`
	Data	interface{}	`json:"data,omniempty"`
	Error	error		`json:"error,omniempty"`
}

func ResponseSuccess(ctx *gin.Context, status int, data interface{}) {
	ctx.JSON(status, Response{
		Status: 	status,
		Message: 	http.StatusText(status),
		Data:		data,
	})
}

func ResponseError(ctx *gin.Context, status int, err error) {
	ctx.JSON(status, Response{
		Status: 	status,
		Message: 	http.StatusText(status),
		Error:		err,
	})
}