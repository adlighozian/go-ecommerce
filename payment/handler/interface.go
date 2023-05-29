package handler

import (
	"github.com/gin-gonic/gin"
)

type Handlerer interface {
	ApprovePayment(ctx *gin.Context)
	CreatePaymentLog(ctx *gin.Context)
}