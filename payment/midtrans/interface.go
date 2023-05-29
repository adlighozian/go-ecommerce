package midtrans

import (
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransInterface interface {
	CheckTransaction(orderID string) (res *coreapi.TransactionStatusResponse, err error)
	CreateTransaction(req *snap.Request) (*snap.Response, error)
}