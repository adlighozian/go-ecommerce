package main

import (
	"product-colors-go/config"
	"product-colors-go/handler"
	"product-colors-go/helper/failerror"
	"product-colors-go/helper/logging"
	"product-colors-go/helper/middleware"
	"product-colors-go/package/db"
	"product-colors-go/publisher"
	"product-colors-go/repository"
	"product-colors-go/service"

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

	repoColors := repository.NewRepository(db.SQLDB, pub)
	Colors := service.NewService(repoColors)
	Handler := handler.NewHandler(Colors)

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
	admin.GET("/colors", hand.GetColors)
	admin.POST("/colors", hand.CreateColors)
	admin.DELETE("/colors", hand.DeleteColors)

	r.Run(conf.Port)
}
