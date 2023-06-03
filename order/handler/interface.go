package handler

import "github.com/gin-gonic/gin"

type Handlerer interface {
	GetOrders(c *gin.Context)
	GetOrdersByStoreID(ctx *gin.Context)
	CreateOrders(c *gin.Context)
	ShowOrders(c *gin.Context)
	UpdateOrders(c *gin.Context)
}
