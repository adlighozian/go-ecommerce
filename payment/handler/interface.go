package handler

import (
	"github.com/gin-gonic/gin"
)

type Handlerer interface {
	GetPaymentMethod(ctx *gin.Context)
	CreatePaymentMethod(ctx *gin.Context)
	CreatePaymentLog(ctx *gin.Context)
}