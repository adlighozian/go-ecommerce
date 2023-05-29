package main

import (
	"log"
	"net/http"
	"payment-go/config"
	"payment-go/handler"
	"payment-go/helper/logging"
	"payment-go/helper/middleware"
	midtransRepo "payment-go/midtrans"
	"payment-go/package/db"
	"payment-go/repository"
	"payment-go/server"
	"payment-go/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"
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

	// inisialization midtrans connection
	// coreapi used to get transaction
	// snapclient used for generating redirect_url to midtrans
	
	var c coreapi.Client
	var s snap.Client
	ServerKey := config.Midtrans.ServerID
	midtrans.ServerKey = ServerKey
	midtrans.Environment = midtrans.Sandbox

	c.New(ServerKey, midtrans.Sandbox)
	s.New(ServerKey, midtrans.Sandbox)

	midtrans := midtransRepo.NewMidtrans(c, s)
	repo := repository.NewRepository(sqlDB.SQLDB, midtrans)
	svc := service.NewService(repo)
	handler := handler.NewHandler(svc)

	// routing
	router := gin.New()
	router.Use(middleware.Logger(logger))
	router.Use(gin.Recovery())

	paymentLog := router.Group("/payments")
	paymentLog.GET("/", handler.CheckTransaction)
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