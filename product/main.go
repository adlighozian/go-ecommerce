package main

import (
	"fmt"
	"product-go/handler"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/middleware"
)

func main() {
	fmt.Println("hallo ini product")

}

func NewServer(hand handler.Handlerer) {
	// server
	r := gin.New()

	// middleware
	r.Use(middleware.Logger())

	r.GET("/products", hand.GetProduct)
	r.GET("/products/details", hand.ShowProduct)
	r.POST("/products", hand.CreateProduct)
	r.PATCH("/products", hand.UpdateProduct)
	r.DELETE("/products", hand.DeleteProduct)

	r.Run(":5000")
}
