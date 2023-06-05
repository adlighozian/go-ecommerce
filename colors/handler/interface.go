package handler

import "github.com/gin-gonic/gin"

type Handlerer interface {
	GetColors(c *gin.Context)
	CreateColors(c *gin.Context)
	DeleteColors(c *gin.Context)
}
