package handler

import (
	"github.com/gin-gonic/gin"
)

type Handlerer interface {
	Get(ctx *gin.Context)
	GetDetail(ctx *gin.Context)
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
}