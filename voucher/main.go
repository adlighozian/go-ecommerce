package main

import (
	"voucher-go/config"
	"voucher-go/handler"
	"voucher-go/helper/failerror"
	"voucher-go/helper/logging"
	"voucher-go/helper/middleware"
	"voucher-go/package/db"
	"voucher-go/publisher"
	"voucher-go/repository"
	"voucher-go/service"

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

	repoVoucher := repository.NewRepository(db.SQLDB, pub)
	Voucher := service.NewService(repoVoucher)
	Handler := handler.NewHandler(Voucher)

	NewServer(Handler, logger)
}

func NewServer(hand handler.Handlerer, logger *zerolog.Logger) {
	// server
	r := gin.New()

	// middleware
	r.Use(middleware.Logger(logger))

	r.GET("/voucher", hand.GetVoucher)
	r.GET("/voucher/details", hand.ShowVoucher)

	admin := r.Group("/admin")
	admin.POST("/voucher", hand.CreateVoucher)
	admin.DELETE("/voucher", hand.DeleteVoucher)

	r.Run(":5002")
}
