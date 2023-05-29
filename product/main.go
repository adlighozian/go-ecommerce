package main

import (
	"product-go/config"
	"product-go/handler"
	"product-go/helper/failerror"
	"product-go/helper/logging"
	"product-go/helper/middleware"
	"product-go/package/db"
	"product-go/publisher"
	"product-go/repository"
	"product-go/service"

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
	// server
	r := gin.New()

	// middleware
	r.Use(middleware.Logger(logger))

	r.GET("/products", hand.GetProduct)
	r.GET("/products/details", hand.ShowProduct)
	
	admin := r.Group("/admin")
	admin.POST("/products", hand.CreateProduct)
	admin.PATCH("/products", hand.UpdateProduct)
	admin.DELETE("/products", hand.DeleteProduct)

	r.Run(":5000")
}
