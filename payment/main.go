package main

import (
	"log"
	"net/http"
	"payment-go/config"
	"payment-go/handler"
	"payment-go/helper/logging"
	"payment-go/helper/middleware"
	"payment-go/midtrans"
	"payment-go/package/db"
	"payment-go/publisher"
	"payment-go/repository"
	"payment-go/server"
	"payment-go/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
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

	var c coreapi.Client
	var s snap.Client
	midtrans := midtrans.NewMidtrans(c, s)
	publisher := publisher.NewPublisher()
	repo := repository.NewRepository(sqlDB.SQLDB, midtrans, publisher)
	svc := service.NewService(repo)
	handler := handler.NewHandler(svc)

	router := gin.New()
	router.Use(middleware.Logger(logger))
	router.Use(gin.Recovery())

	paymentLog := router.Group("/payments")
	paymentLog.GET("/", handler.ApprovePayment)
	paymentLog.POST("/", handler.CreatePaymentLog)

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