package main

import (
	"api-gateway-go/config"
	"api-gateway-go/handler"
	"api-gateway-go/helper/logging"
	"api-gateway-go/helper/middleware"
	"api-gateway-go/package/db"
	"api-gateway-go/repository"
	"api-gateway-go/server"
	"api-gateway-go/service"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
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

	pingHandler := handler.NewPingGinHandler()

	shortenRepo := repository.NewShortenRepo(sqlDB.SQLDB)
	shortenSvc := service.NewShortenService(shortenRepo)
	shortenHandler := handler.NewShortenHandler(shortenSvc)

	router := gin.New()
	router.Use(cors.Default())
	router.Use(requestid.New())
	router.Use(middleware.Logger(logger))
	router.Use(gin.Recovery())
	// if config.Debug {
	// 	pprof.Register(router)
	// }

	router.GET("/ping", pingHandler.Ping)

	router.POST("/shorten", shortenHandler.Shorten)
	router.Any("/:hash", shortenHandler.Get)

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
