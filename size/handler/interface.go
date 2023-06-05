package handler

import "github.com/gin-gonic/gin"

type Handlerer interface {
	GetSize(c *gin.Context)
	CreateSize(c *gin.Context)
	DeleteSize(c *gin.Context)
}
