package response

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type ResSuccess struct {
	Status	int			`json:"status"`
	Message string		`json:"message"`
	Data	interface{}	`json:"data"`
}

type ResError struct {
	Status	int			`json:"status"`
	Message string		`json:"message"`
	Error	string		`json:"error"`
}

func ResponseSuccess(ctx *gin.Context, status int, data interface{}) {
	ctx.JSON(status, ResSuccess{
		Status: 	status,
		Message: 	http.StatusText(status),
		Data:		data,
	})
}

func ResponseError(ctx *gin.Context, status int, err error) {
	ctx.JSON(status, ResError{
		Status: 	status,
		Message: 	http.StatusText(status),
		Error:		err.Error(),
	})
}