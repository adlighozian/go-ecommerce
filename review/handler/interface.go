package handler

import (
	"github.com/gin-gonic/gin"
)

type Handlerer interface {
	GetByProductID(ctx *gin.Context)
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
}