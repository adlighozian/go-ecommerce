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

func ResponseSuccess(ctx *gin.Context, status int, message string, data interface{}) {
	if message == "" {
		message = http.StatusText(status)
	}
	ctx.JSON(status, Response{
		Status: 	status,
		Message: 	message,
		Data:		data,
	})
}

func ResponseError(ctx *gin.Context, status int, message string, err error) {
	if message == "" {
		message = http.StatusText(status)
	}
	ctx.JSON(status, Response{
		Status: 	status,
		Message: 	message,
		Error:		err,
	})
}