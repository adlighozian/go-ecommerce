package repository

import (
	"consumer-payment-logs-go/model"
)

type Repositorier interface {
	Create(req model.PaymentLogRequest) (err error)
}