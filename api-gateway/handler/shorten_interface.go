package handler

import "github.com/gin-gonic/gin"

type ShortenHandlerI interface {
	Get(ctx *gin.Context)
	Shorten(ctx *gin.Context)
}
