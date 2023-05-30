package handler

import (
	"encoding/json"
	"user-consumer-go/model"
	"user-consumer-go/package/rmq"
	"user-consumer-go/service"

	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type UserSettingHandler struct {
	rmq    rmq.RabbitMQClient
	logger *zerolog.Logger
	svc    service.UserSettingServiceI
}

func NewUserSettingHandler(
	rmq rmq.RabbitMQClient, logger *zerolog.Logger, svc service.UserSettingServiceI,
) UserSettingHandlerI {
	h := new(UserSettingHandler)
	h.rmq = rmq
	h.logger = logger
	h.svc = svc
	return h
}

func (h *UserSettingHandler) UpdateByUserID() {
	msgs, errSub := h.rmq.Subscribe(
		"user_setting.updated",
		"topic",
		"user_setting.updated",
		"user_setting.updated",
		false,
		1,
	)
	if errSub != nil {
		h.logger.Error().Err(errSub).Msg("rmq.Subscribe err")
		return
	}

	go func(deliveries <-chan amqp091.Delivery) {
		newSetting := new(model.UserSetting)
		for d := range deliveries {
			_ = json.Unmarshal(d.Body, &newSetting)

			// h.logger.Debug().Msgf("%v", newSetting)

			setting, ucErr := h.svc.UpdateByUserID(newSetting)
			if ucErr != nil {
				h.logger.Error().Err(ucErr).Msgf("svc.UpdateByID err: %v\n%v", newSetting, setting)
				// _ = d.Nack(false, false)
				continue
			}

			_ = d.Ack(false)
			h.logger.Debug().Msgf("from %s, rmq.Sub success with id:%v", d.RoutingKey, setting.ID)
		}
	}(msgs)
}
