package main

import (
	"github.com/gin-gonic/gin"
	"wishlist-go/package/db"
	"wishlist-go/config"
	"wishlist-go/handler"
	"wishlist-go/service"
	"wishlist-go/repository"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	db, err := db.NewGormDB(true, "pgx", conf.DatabaseURL)
	if err != nil {
		panic(err)
	}

	repo := repository.NewRepository(db.SqlDB)
	svc := service.NewService(repo)
	handler := handler.NewHandler(svc)

	r := gin.Default()
	review := r.Group("/wishlist")
	review.GET("/", handler.GetDetail)
	review.POST("/", handler.Create)
	review.DELETE("/", handler.Delete)

	r.Run(":5000")
}