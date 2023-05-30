package main

import (
	"encoding/json"
	"log"
	amqp "github.com/rabbitmq/amqp091-go"

	"consumer-wishlist-go/config"
	"consumer-wishlist-go/package/db"
	"consumer-wishlist-go/helper/logging"
	"consumer-wishlist-go/model"
	"consumer-wishlist-go/repository"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		FailOnError(err, "failed to load config")
		return
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

	repo := repository.NewRepository(sqlDB.SQLDB)

	conn, err := amqp.Dial(config.RabbitMQURL)	
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"create_wishlists", 	  // queue name
		true,                 // durable
		false,                // auto delete queue when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue name
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	FailOnError(err, "Failed to register a consumer")

	// consumer must always be on and the channel to prevent the consumer from turning off
	var forever chan string

	// worker to receive value from variable msgs
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			var data []model.WishlistRequest
			err := json.Unmarshal(d.Body, &data)
			if err != nil {
				FailOnError(err, "error unmarshal")
			}

			err = repo.Create(data)
			if err != nil {
				FailOnError(err, "error create bulk")
			}
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	// channel in to prevent consumer to turning off
	<-forever
}