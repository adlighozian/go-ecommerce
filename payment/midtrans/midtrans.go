package midtrans

import (
	"fmt"
	"log"
	"payment-go/config"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

var s snap.Client
var c coreapi.Client

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

func (m Midtrans) ApprovePayment(orderID string) (res *coreapi.ChargeResponse, err error) {
	// inisialization snap connection
	config, confErr := config.LoadConfig()
	if confErr != nil {
		log.Fatalf("load config err:%s", confErr)
	}

	ServerKey := config.Midtrans.ServerID
	midtrans.ServerKey = ServerKey
	midtrans.Environment = midtrans.Sandbox

	c.New(ServerKey, midtrans.Sandbox)
	
	// approve transaction
	res, err = c.ApproveTransaction(orderID)
	if res != nil {
		return
	}
	return
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

	s.New(ServerKey, midtrans.Sandbox)

	// create transaction
	resp, err := s.CreateTransaction(req)
	if err != nil {
		return nil, fmt.Errorf(err.GetMessage())
	}
	fmt.Println("RESPONSE 1: ", resp)
	return resp, nil
}