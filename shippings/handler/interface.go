package handler

import "github.com/gin-gonic/gin"

type Handlerer interface {
	GetShipping(c *gin.Context)
	CreateShipping(c *gin.Context)
	DeleteShipping(c *gin.Context)
}
