package handler

import "github.com/gin-gonic/gin"

type UserHandlerI interface {
	GetByID(ctx *gin.Context)
	UpdateByID(ctx *gin.Context)
}
