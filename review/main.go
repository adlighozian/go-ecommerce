package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"review-go/package/db"
	"review-go/config"
	"review-go/handler"
	"review-go/helper/logging"
	"review-go/helper/middleware"
	"review-go/publisher"
	"review-go/service"
	"review-go/server"
	"review-go/repository"
	"time"
)

func main() {
	config, confErr := config.LoadConfig()
	if confErr != nil {
		log.Fatalf("load config err:%s", confErr)
	}

	logger := logging.New(config.Debug)

	sqlDB, errDB := db.NewGormDB(config.Debug, config.Database.Driver, config.Database.URL)
	if errDB != nil {
		logger.Fatal().Err(errDB).Msg("db failed to connect")
	}
	logger.Debug().Msg("db connected")

	defer func() {
		logger.Debug().Msg("closing db")
		_ = sqlDB.Close()
	}()

	publisher := publisher.NewPublisher()
	repo := repository.NewRepository(sqlDB.SQLDB, publisher)
	svc := service.NewService(repo)
	handler := handler.NewHandler(svc)

	router := gin.New()
	router.Use(middleware.Logger(logger))
	router.Use(gin.Recovery())

	review := router.Group("/reviews")
	review.GET("/", handler.GetByProductID)
	review.POST("/", handler.Create)

	srv := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if srvErr := server.Run(srv, logger); srvErr != nil {
		logger.Fatal().Err(srvErr).Msg("server shutdown failed")
	}
}