package handler

import (
	"github.com/gin-gonic/gin"
)

type Handlerer interface {
	Get(c *gin.Context)
	GetDetail(c *gin.Context)
	Create(c *gin.Context)
	Delete(c *gin.Context)
}