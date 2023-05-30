package handler

import "github.com/gin-gonic/gin"

type AuthHandlerI interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)

	GoogleLogin(ctx *gin.Context)
	GoogleCallback(ctx *gin.Context)
}
