package main

import (
	"cart-go/config"
	"cart-go/handler"
	"cart-go/helper/logging"
	"cart-go/helper/middleware"
	"cart-go/package/db"
	"cart-go/publisher"
	"cart-go/repository"
	"cart-go/server"
	"cart-go/service"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

	review := router.Group("/carts")
	review.GET("/", handler.Get)
	review.POST("/", handler.Create)
	review.DELETE("/", handler.Delete)

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