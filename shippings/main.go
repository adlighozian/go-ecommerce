package main

import (
	"shippings-go/config"
	"shippings-go/handler"
	"shippings-go/helper/failerror"
	"shippings-go/helper/logging"
	"shippings-go/helper/middleware"
	"shippings-go/package/db"
	"shippings-go/publisher"
	"shippings-go/repository"
	"shippings-go/service"

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

	repoShipping := repository.NewRepository(db.SQLDB, pub)
	Shipping := service.NewService(repoShipping)
	Handler := handler.NewHandler(Shipping)

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
	admin.GET("/shipping", hand.GetShipping)
	admin.POST("/shipping", hand.CreateShipping)
	admin.DELETE("/shipping", hand.DeleteShipping)

	r.Run(conf.Port)
}
