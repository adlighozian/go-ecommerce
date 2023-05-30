package midtrans

import (
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransInterface interface {
	CheckPayment(orderID string) (*coreapi.TransactionStatusResponse, error)
	CreatePayment(req *snap.Request) (*snap.Response, error)
}