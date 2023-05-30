package midtrans

import (
	"fmt"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

// coreapi used to get transaction
// snapclient used for generating redirect_url to midtrans
type Midtrans struct {
	coreapiclient coreapi.Client
	snapclient 	  snap.Client
}

func NewMidtrans(coreapiclient coreapi.Client, snapclient snap.Client) MidtransInterface {
	return Midtrans{
		coreapiclient: coreapiclient,
		snapclient: 	snapclient,
	}
}

func (m Midtrans) CheckPayment(orderID string) (*coreapi.TransactionStatusResponse, error) {
	// get transaction status by order id in midtrans
	resp, err := m.coreapiclient.CheckTransaction(orderID)
	if err != nil {
		return nil, fmt.Errorf(err.GetMessage())
	}
	return resp, nil
}

func(m Midtrans) CreatePayment(req *snap.Request) (*snap.Response, error) {
	// create transaction in midtrans
	resp, err := m.snapclient.CreateTransaction(req)
	if err != nil {
		return nil, fmt.Errorf(err.GetMessage())
	}

	return resp, nil
}