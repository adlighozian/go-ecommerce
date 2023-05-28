package service

import (
	"payment-go/model"
)

type Servicer interface {
	GetPaymentMethod() (res []model.PaymentMethod, err error)
	GetPaymentMethodByName(name string) (res model.PaymentMethod, err error)
	CreatePaymentMethod(req []model.PaymentMethodRequest) (res []model.PaymentMethod, err error)
	CreatePaymentLog(req model.PaymentLogRequest) (res model.PaymentLog, err error)
	DeletePaymentMethod(paymentMethodID int) (err error)
}