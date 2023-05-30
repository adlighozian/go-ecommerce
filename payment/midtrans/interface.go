package midtrans

import (
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransInterface interface {
	CheckTransaction(orderID string) (*coreapi.TransactionStatusResponse, error)
	CreateTransaction(req *snap.Request) (*snap.Response, error)
}