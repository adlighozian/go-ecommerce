package handler

import (
	"encoding/json"
	"user-consumer-go/model"
	"user-consumer-go/package/rmq"
	"user-consumer-go/service"

	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type UserHandler struct {
	rmq    rmq.RabbitMQClient
	logger *zerolog.Logger
	svc    service.UserServiceI
}

func NewUserHandler(rmq rmq.RabbitMQClient, logger *zerolog.Logger, svc service.UserServiceI) UserHandlerI {
	h := new(UserHandler)
	h.rmq = rmq
	h.logger = logger
	h.svc = svc
	return h
}

func (h *UserHandler) Create() {
	msgs, errSub := h.rmq.Subscribe(
		"user.created",
		"topic",
		"user.created",
		"user.created",
		false,
		1,
	)
	if errSub != nil {
		h.logger.Error().Err(errSub).Msg("rmq.Subscribe err")
		return
	}

	go func(deliveries <-chan amqp091.Delivery) {
		newUser := new(model.User)
		for d := range deliveries {
			_ = json.Unmarshal(d.Body, &newUser)

			// h.logger.Debug().Msgf("%v", newUser)

			user, ucErr := h.svc.Create(newUser)
			if ucErr != nil {
				h.logger.Error().Err(ucErr).Msg("svc.Create err")
				// _ = d.Nack(false, false)
				continue
			}

			_ = d.Ack(false)
			h.logger.Debug().Msgf("from %s, rmq.Sub success with id:%v", d.RoutingKey, user.ID)
		}
	}(msgs)
}
