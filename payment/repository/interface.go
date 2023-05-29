package repository

import (
	"payment-go/model"

	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type Repositorier interface {
	ApprovePayment(orderID string) (res *coreapi.ChargeResponse, err error)
	CreatePaymentLog(req model.PaymentLogRequest) (res *snap.Response, err error)
}
