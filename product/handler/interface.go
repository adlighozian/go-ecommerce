package handler

import "github.com/gin-gonic/gin"

type Handlerer interface {
	GetProduct(c *gin.Context)
	ShowDetail(c *gin.Context)
	CreateProduct(c *gin.Context)
	UpdateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
}
