package main

import (
	"consumer-insert-shipping-go/config"
	"consumer-insert-shipping-go/db"
	"consumer-insert-shipping-go/helpers"
	"consumer-insert-shipping-go/model"
	"consumer-insert-shipping-go/repository"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conf, err := config.LoadConfig()
	helpers.FailOnError(err, "error loadconfig")

	conn, err := amqp.Dial(conf.RabbitMQ)
	helpers.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	helpers.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"create_shipping", // name
		true,              // durable
		false,             // auto delete queue when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	helpers.FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	helpers.FailOnError(err, "Failed to register a consumer")

	// consumer must always be on and the channel to prevent the consumer from turning off
	var forever chan string

	db, err := db.NewGormDB(conf.Debug, conf.DatabaseDriver, conf.DatabaseURL)
	helpers.FailOnError(err, "error loadconfig")

	// worker to receive value from variable msgs
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			var data []model.Shipping
			err := json.Unmarshal(d.Body, &data)
			if err != nil {
				helpers.FailOnError(err, "error unmarshal")
			}

			err = repository.NewProduct(db.SQLDB).CreateProduct(data)
			if err != nil {
				helpers.FailOnError(err, "error create bulk transaction detail")
			}
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	// channel in to prevent consumer to turning off
	<-forever
}
