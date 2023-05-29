package handler

import (
	"github.com/gin-gonic/gin"
)

type Handlerer interface {
	CheckTransaction(ctx *gin.Context)
	CreatePaymentLog(ctx *gin.Context)
}