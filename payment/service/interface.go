package service

import (
	"payment-go/model"
)

type Servicer interface {
	GetPaymentMethod() (res []model.PaymentMethod, err error)
	CreatePaymentMethod(req []model.PaymentMethodRequest) (res []model.PaymentMethod, err error)
	CreatePaymentLog(req []model.PaymentLogsRequest) (res []model.PaymentLogsRequest, err error)
}