package main

import (
	"log"
	"net/http"
	"splash-screen-go/config"
	"splash-screen-go/handler"
	"splash-screen-go/helper/logging"
	"splash-screen-go/helper/middleware"
	"splash-screen-go/package/db"
	"splash-screen-go/publisher"
	"splash-screen-go/repository"
	"splash-screen-go/server"
	"splash-screen-go/service"
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

	review := router.Group("/splashscreens")
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
