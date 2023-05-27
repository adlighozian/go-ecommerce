package main

import (
	"api-gateway-go/config"
	"api-gateway-go/handler"
	"api-gateway-go/helper/logging"
	"api-gateway-go/helper/middleware"
	"api-gateway-go/package/db"
	"api-gateway-go/package/redisclient"
	"api-gateway-go/repository"
	"api-gateway-go/server"
	"api-gateway-go/service"
	"log"
	"net/http"
	"time"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
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

	redisClient, errRedis := redisclient.NewRedisClient(config.Redis.Addr, config.Redis.Password, config.Redis.DB)
	if errDB != nil {
		logger.Fatal().Err(errRedis).Msg("redis failed to connect")
	}
	logger.Debug().Msg("redis connected")

	// Initialize Gorm adapter and the Casbin enforcer with the model
	adapter, _ := gormadapter.NewAdapter(config.Database.Driver, config.Database.URL, config.Database.DBName, true)
	enforcer, errEnforcer := casbin.NewEnforcer("./model.conf", adapter)
	if errDB != nil {
		logger.Fatal().Err(errEnforcer).Msg("enforcer failed")
	}
	logger.Debug().Msg("enforcer connected")

	defer func() {
		logger.Debug().Msg("closing db")
		_ = sqlDB.Close()
		logger.Debug().Msg("db closed")

		logger.Debug().Msg("closing redis")
		_ = redisClient.Close()
		logger.Debug().Msg("redis closed")
	}()

	pingHandler := handler.NewPingGinHandler()

	shortenRepo := repository.NewShortenRepo(sqlDB.SQLDB, redisClient.Redis)
	shortenSvc := service.NewShortenService(shortenRepo)
	shortenHandler := handler.NewShortenHandler(shortenSvc)

	router := gin.New()
	router.Use(cors.Default())

	// 9. Error and panic handling
	router.Use(middleware.Logger(logger))
	router.Use(gin.Recovery())

	// 1. Paramater validation,
	// 6. Dynamic routing using path parameters,
	// 7. Service discovery using database,
	router.Use(middleware.HashedURLConverter(shortenSvc))

	// 2. Allow-path
	allowedPaths := []string{
		// "ping",
		"register",
		"login", "login/cms",
	}
	// 3. Authentication
	router.Use(middleware.AuthMiddleware(config.JWTSecretKey, allowedPaths))
	// 4. Authorization
	router.Use(middleware.AuthzMiddleware(enforcer, allowedPaths))

	router.Use(requestid.New())
	// 5. Request counter
	router.Use(middleware.RequestCounter(redisClient.Redis))

	// things outside my capabilities:
	// 2. Allow-list based on IPs
	// 10. Circuit Breaker
	// 11. Monitoring
	// 12. Cache

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
