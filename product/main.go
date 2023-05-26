package main

import (
	"product-go/handler"
	"product-go/package/db"
	"product-go/repository"
	"product-go/service"

	"github.com/gin-gonic/gin"
)

func main() {
	db := db.GetConnection()

	repoProduct := repository.NewRepository(db)
	product := service.NewService(repoProduct)
	Handler := handler.NewHandler(product)

	NewServer(Handler)
}

func NewServer(hand handler.Handlerer) {
	// server
	r := gin.New()

	// middleware
	// r.Use(middleware.Logger())

	r.GET("/products", hand.GetProduct)
	r.GET("/products/details", hand.ShowProduct)
	r.POST("/products", hand.CreateProduct)
	r.PATCH("/products", hand.UpdateProduct)
	r.DELETE("/products", hand.DeleteProduct)

	r.Run(":5000")
}
