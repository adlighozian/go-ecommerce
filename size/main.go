package main

import (
	"size-go/config"
	"size-go/handler"
	"size-go/helper/failerror"
	"size-go/helper/logging"
	"size-go/helper/middleware"
	"size-go/package/db"
	"size-go/publisher"
	"size-go/repository"
	"size-go/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func main() {

	conf, err := config.LoadConfig()
	failerror.FailError(err, "error loadconfig")

	logger := logging.New(conf.Debug)

	db, err := db.NewGormDB(conf.Debug, conf.DatabaseDriver, conf.DatabaseURL)
	failerror.FailError(err, "error new gorm")
	logger.Debug().Msg("DB Connected")

	pub := publisher.NewPublisher()

	repoSize := repository.NewRepository(db.SQLDB, pub)
	Size := service.NewService(repoSize)
	Handler := handler.NewHandler(Size)

	NewServer(Handler, logger)
}

func NewServer(hand handler.Handlerer, logger *zerolog.Logger) {
	conf, err := config.LoadConfig()
	failerror.FailError(err, "error loadconfig")
	// server
	r := gin.New()

	// middleware
	r.Use(middleware.Logger(logger))

	admin := r.Group("/admin")
	admin.GET("/size", hand.GetSize)
	admin.POST("/size", hand.CreateSize)
	admin.DELETE("/size", hand.DeleteSize)

	r.Run(conf.Port)
}
