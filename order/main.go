package main

import (
	"order-go/config"
	"order-go/handler"
	"order-go/helper/failerror"
	"order-go/helper/logging"
	"order-go/helper/middleware"
	"order-go/package/db"
	"order-go/publisher"
	"order-go/repository"
	"order-go/service"

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

	repoProduct := repository.NewRepository(db.SQLDB, pub)
	product := service.NewService(repoProduct)
	Handler := handler.NewHandler(product)

	NewServer(Handler, logger)
}

func NewServer(hand handler.Handlerer, logger *zerolog.Logger) {
	conf, err := config.LoadConfig()
	failerror.FailError(err, "error loadconfig")

	// server
	r := gin.New()

	// middleware
	r.Use(middleware.Logger(logger))

	r.GET("/orders", hand.GetOrders)
	r.GET("/orders/stores", hand.GetOrdersByStoreID)
	r.GET("/orders/details", hand.ShowOrders)
	r.POST("/orders", hand.CreateOrders)
	r.POST("/orders/item", hand.CreateOrders)

	admin := r.Group("/admin")
	admin.PATCH("/orders", hand.UpdateOrders)

	r.Run(":" + conf.Port)
}
