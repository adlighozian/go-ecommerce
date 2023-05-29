package midtrans

import (
	"fmt"
	"log"
	"payment-go/config"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

// coreapi used for approve transaction
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

func (m Midtrans) CheckTransaction(orderID string) (*coreapi.TransactionStatusResponse, error) {
	// inisialization snap connection
	config, confErr := config.LoadConfig()
	if confErr != nil {
		log.Fatalf("load config err:%s", confErr)
	}

	ServerKey := config.Midtrans.ServerID
	midtrans.ServerKey = ServerKey
	midtrans.Environment = midtrans.Sandbox

	m.coreapiclient.New(ServerKey, midtrans.Sandbox)
	
	// get transaction status by order id
	resp, err := m.coreapiclient.CheckTransaction(orderID)
	if err != nil {
		return nil, fmt.Errorf(err.GetMessage())
	}
	return resp, nil
}

func(m Midtrans) CreateTransaction(req *snap.Request) (*snap.Response, error) {
	// inisialization snap connection
	config, confErr := config.LoadConfig()
	if confErr != nil {
		log.Fatalf("load config err:%s", confErr)
	}

	ServerKey := config.Midtrans.ServerID
	midtrans.ServerKey = ServerKey
	midtrans.Environment = midtrans.Sandbox

	m.snapclient.New(ServerKey, midtrans.Sandbox)

	// create transaction
	resp, err := m.snapclient.CreateTransaction(req)
	if err != nil {
		return nil, fmt.Errorf(err.GetMessage())
	}

	return resp, nil
}