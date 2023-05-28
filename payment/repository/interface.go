package repository

import (
	"payment-go/model"
)

type Repositorier interface {
	GetPaymentMethod() (res []model.PaymentMethod, err error)
	GetPaymentMethodByID(paymentMethodID int) (res model.PaymentMethod, err error)
	GetPaymentMethodByName(name string) (res model.PaymentMethod, err error)
	CreatePaymentMethod(req []model.PaymentMethodRequest) (err error)
	CreatePaymentLog(req model.PaymentLogRequest) (res model.PaymentLog, err error)
	DeletePaymentMethod(paymentMethodID int) (err error)
}