package repository

import (
	"payment-go/model"

	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type Repositorier interface {
	CheckTransaction(orderID string) (res *coreapi.TransactionStatusResponse, err error)
	CreatePaymentLog(req model.PaymentLogRequest) (res *snap.Response, err error)
}