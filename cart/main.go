package main

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"cart-go/package/db"
	"cart-go/config"
	"cart-go/handler"
	"cart-go/helper/logging"
	"cart-go/helper/middleware"
	"cart-go/service"
	"cart-go/server"
	"cart-go/repository"
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

	repo := repository.NewRepository(sqlDB.SQLDB)
	svc := service.NewService(repo)
	handler := handler.NewHandler(svc)

	router := gin.New()
	router.Use(middleware.Logger(logger))
	router.Use(gin.Recovery())

	review := router.Group("/cart")
	review.GET("/", handler.GetDetail)
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