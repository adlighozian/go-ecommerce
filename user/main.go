package main

import (
	"log"
	"net/http"
	"time"
	"user-go/config"
	"user-go/handler"
	"user-go/helper/logging"
	"user-go/helper/middleware"
	"user-go/package/db"
	"user-go/package/redisclient"
	"user-go/package/rmq"
	"user-go/repository"
	"user-go/server"
	"user-go/service"

	"github.com/gin-gonic/gin"
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

	redisClient, errRedis := redisclient.NewRedisClient(
		config.Redis.Addr, config.Redis.ClientName,
		config.Redis.Username, config.Redis.Password,
		config.Redis.DB,
	)
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

	userRepo := repository.NewUserRepository(sqlDB.SQLDB, redisClient.Redis, rmq)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	userSettingRepo := repository.NewUserSettingRepository(sqlDB.SQLDB, redisClient.Redis, rmq)
	userSettingSvc := service.NewUserSettingService(userSettingRepo)
	userSettingHandler := handler.NewUserSettingHandler(userSettingSvc)

	pingHandler := handler.NewPingGinHandler()

	router := gin.New()
	router.Use(middleware.Logger(logger))
	router.Use(gin.Recovery())

	userRouter := router.Group("/user")
	{
		userRouter.GET("/ping", pingHandler.Ping)

		userRouter.GET("/profiles", userHandler.GetByID)
		userRouter.PATCH("/profiles", userHandler.UpdateByID)

		userRouter.PATCH("/profiles/settings", userSettingHandler.UpdateByUserID)
	}

	srv := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Debug().Msgf("service will be start at port: %v", config.Port)

	if errSrv := server.Run(srv, logger); errSrv != nil {
		logger.Fatal().Err(errSrv).Msg("server shutdown failed")
	}
}
