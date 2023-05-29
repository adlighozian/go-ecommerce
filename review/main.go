package main

import (
	"github.com/gin-gonic/gin"
	"review-go/package/db"
	"review-go/config"
	"review-go/handler"
	"review-go/service"
	"review-go/repository"
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
	review := r.Group("/review")
	review.GET("/", handler.GetByProductID)
	review.POST("/", handler.Create)

	r.Run(":5000")
}