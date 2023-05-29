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

func (m Midtrans) CheckTransaction(orderID string) (res *coreapi.TransactionStatusResponse, err error) {
	// inisialization snap connection
	config, confErr := config.LoadConfig()
	if confErr != nil {
		log.Fatalf("load config err:%s", confErr)
	}

	ServerKey := config.Midtrans.ServerID
	midtrans.ServerKey = ServerKey
	midtrans.Environment = midtrans.Sandbox

	c.New(ServerKey, midtrans.Sandbox)
	
	// get transaction status by order id
	res, err = c.CheckTransaction(orderID)
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

	fmt.Println("ORDER ID : ", req.TransactionDetails.OrderID)
	c.New(ServerKey, midtrans.Sandbox)
	res, _ := c.CheckTransaction(req.TransactionDetails.OrderID)
	fmt.Println("CHECK TRANSACTION: ", res)
	fmt.Println("CHECK TRANSACTION ORDER ID : ", res.OrderID)

	return resp, nil
}