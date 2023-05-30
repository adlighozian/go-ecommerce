package publisher

import (
	"encoding/json"
	"log"
	"order-go/config"
	"order-go/helper/failerror"
	"order-go/helper/timeout"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher interface {
	Public(req any, queueName string) error
}

type publisher struct{}

func NewPublisher() Publisher {
	return &publisher{}
}

func (p publisher) Public(req any, queueName string) error {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	config, err := config.LoadConfig()
	failerror.FailError(err, "error config")

	conn, err := amqp.Dial(config.RabbitMQ)
	failerror.FailError(err, "error connect to rabbitmq")
	defer conn.Close()

	ch, err := conn.Channel()
	failerror.FailError(err, "failed to open a channel")

	q, err := ch.QueueDeclare(
		queueName, // queue name
		true,      // durable
		false,     // auto delete queue when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failerror.FailError(err, "failed to declare queue")

	// marshal data to jsonByte
	jsonByte, err := json.Marshal(req)
	failerror.FailError(err, "failed to marshaling")

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         jsonByte,
		})
	failerror.FailError(err, "failed to publish message")

	log.Printf(" [x] Sent %s", req)

	return nil
}
