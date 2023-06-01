package handler

import "github.com/gin-gonic/gin"

type Handlerer interface {
	GetVoucher(c *gin.Context)
	ShowVoucher(c *gin.Context)
	CreateVoucher(c *gin.Context)
	DeleteVoucher(c *gin.Context)
}
