package main

import (
	"auth-go/config"
	"auth-go/handler"
	"auth-go/helper/logging"
	"auth-go/helper/middleware"
	"auth-go/package/db"
	"auth-go/package/redisclient"
	"auth-go/package/rmq"
	"auth-go/repository"
	"auth-go/server"
	"auth-go/service"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	config, errConf := config.LoadConfig()
	if errConf != nil {
		log.Fatalf("load config err:%s", errConf)
	}

	logger := logging.New(config.Debug)

	sqlDB, errDB := db.NewGormDB(config.Debug, config.Database.Driver, config.Database.URL)
	if errDB != nil {
		logger.Fatal().Err(errDB).Msg("db failed to connect")
	}
	logger.Debug().Msg("db connected")

	redisClient, errRedis := redisclient.NewRedisClient(config.Redis.Addr, config.Redis.Password, config.Redis.DB)
	if errDB != nil {
		logger.Fatal().Err(errRedis).Msg("redis failed to connect")
	}
	logger.Debug().Msg("redis connected")

	rmq, errRmq := rmq.NewRabbitMQ(config.RabbitMQ.URL)
	if errRmq != nil {
		logger.Fatal().Err(errRmq).Msg("rabbitmq failed to connect")
	}
	logger.Debug().Msg("rabbitmq connected")

	defer func() {
		errDBC := sqlDB.Close()
		if errDBC != nil {
			logger.Fatal().Err(errDBC).Msg("db failed to closed")
		}
		logger.Debug().Msg("db closed")

		errRedisC := redisClient.Close()
		if errRedisC != nil {
			logger.Fatal().Err(errRedisC).Msg("redis failed to closed")
		}
		logger.Debug().Msg("redis closed")

		errRmqC := rmq.Shutdown()
		if errRmqC != nil {
			logger.Fatal().Err(errRmqC).Msg("rabbitmq failed to closed")
		}
		logger.Debug().Msg("rabbitmq closed")
	}()

	gauth := &oauth2.Config{
		ClientID:     config.GoogleClientID,
		ClientSecret: config.GoogleClientSecret,
		RedirectURL:  "http://localhost:8081/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	authRepo := repository.NewAuthRepository(sqlDB.SQLDB, redisClient.Redis, rmq)
	authSvc := service.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authSvc, config.JWTSecretKey, gauth)

	pingHandler := handler.NewPingGinHandler()

	router := gin.New()
	router.Use(middleware.Logger(logger))
	router.Use(gin.Recovery())

	authRouter := router.Group("/auth")
	{
		authRouter.GET("/ping", pingHandler.Ping)

		authRouter.POST("/register", authHandler.Register)
		authRouter.POST("/login", authHandler.Login)
		authRouter.POST("/refresh-token", authHandler.RefreshToken)

		authRouter.GET("/google/login", authHandler.GoogleLogin)
		authRouter.GET("/google/callback", authHandler.GoogleCallback)
	}

	srv := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if errSrv := server.Run(srv, logger); errSrv != nil {
		logger.Fatal().Err(errSrv).Msg("server shutdown failed")
	}
}
