package handler

import (
	"github.com/gin-gonic/gin"
)

type Handlerer interface {
	CheckPayment(ctx *gin.Context)
	CreatePaymentLog(ctx *gin.Context)
}
