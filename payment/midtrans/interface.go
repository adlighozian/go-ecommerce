package midtrans

import (
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransInterface interface {
	ApprovePayment(orderID string) (res *coreapi.ChargeResponse, err error)
	CreateTransaction(req *snap.Request) (*snap.Response, error)
}