package rmq

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/rabbitmq/amqp091-go"
)

type Response struct {
	sync.RWMutex
	Map map[string]chan string
}

type RabbitMQClient interface {
	Publish(ctx context.Context, exchangeName, exchangeType, contentType, routingKeyPub string, dataBytes []byte) error
	Subscribe(
		exchangeName, exchangeType, queueName, routingKeyCon string,
		autoAck bool,
		prefetchCount int,
	) (<-chan amqp091.Delivery, error)
	ConnClose() error
	ChClose() error
	Shutdown() error
}

type RabbitMQ struct {
	url  string
	Conn *amqp091.Connection
	Ch   *amqp091.Channel
	done chan bool
}

func NewRabbitMQ(url string) (RabbitMQClient, error) {
	if url == "" {
		return nil, errors.New("no rabbitmq url")
	}

	rmq := new(RabbitMQ)
	rmq.url = url
	rmq.done = make(chan bool)

	err := rmq.init()
	if err != nil {
		return nil, err
	}

	return rmq, nil
}

func (r *RabbitMQ) init() error {
	conn, errConn := amqp091.Dial(r.url)
	if errConn != nil {
		return fmt.Errorf("failed to connect to rabbitmq: %w", errConn)
	}

	ch, errCh := conn.Channel()
	if errCh != nil {
		return fmt.Errorf("failed to open a channel: %w", errCh)
	}

	r.Conn = conn
	r.Ch = ch
	return nil
}

func (r *RabbitMQ) Publish(
	ctx context.Context,
	exchangeName, exchangeType, contentType, routingKeyPub string,
	dataBytes []byte,
) error {
	errExc := r.Ch.ExchangeDeclare(
		exchangeName, // exchange name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if errExc != nil {
		return errExc
	}

	errPub := r.Ch.PublishWithContext(
		ctx,
		exchangeName,  // exchange name
		routingKeyPub, // routing key
		false,         // mandatory
		false,         // immediate
		amqp091.Publishing{
			ContentType: contentType,
			Body:        dataBytes,
		})
	if errPub != nil {
		return errPub
	}

	return nil
}

func (r *RabbitMQ) Subscribe(
	exchangeName, exchangeType, queueName, routingKeyCon string,
	autoAck bool,
	prefetchCount int,
) (<-chan amqp091.Delivery, error) {
	errExc := r.Ch.ExchangeDeclare(
		exchangeName, // exchange name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if errExc != nil {
		return nil, errExc
	}

	q, errQ := r.Ch.QueueDeclare(
		queueName, // queue name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if errQ != nil {
		return nil, errQ
	}

	errQB := r.Ch.QueueBind(
		q.Name,        // queue name
		routingKeyCon, // routing key
		exchangeName,  // exchange
		false,         // no-wait
		nil,           // arguments
	)
	if errQB != nil {
		return nil, errQB
	}

	errQos := r.Ch.Qos(
		prefetchCount, // prefetch count
		0,             // prefetch size
		false,         // global
	)
	if errQos != nil {
		return nil, errQos
	}

	msgs, errCon := r.Ch.Consume(
		q.Name,  // queue name
		"",      // consumer tag
		autoAck, // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if errCon != nil {
		return nil, errCon
	}
	return msgs, nil
}

func (r *RabbitMQ) ConnClose() error {
	return r.Conn.Close()
}

func (r *RabbitMQ) ChClose() error {
	return r.Ch.Close()
}

func (r *RabbitMQ) Shutdown() error {
	if err := r.ChClose(); err != nil {
		return err
	}
	if err := r.ConnClose(); err != nil {
		return err
	}

	close(r.done)

	return nil
}
