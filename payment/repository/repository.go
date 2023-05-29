package repository

import (
	"database/sql"
	"fmt"
	midtransrepo "payment-go/midtrans"
	"payment-go/model"
	"payment-go/publisher"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type repository struct {
	db        *sql.DB
	midtrans  midtransrepo.MidtransInterface
	publisher publisher.PublisherInterface
}

func NewRepository(db *sql.DB, midtrans midtransrepo.MidtransInterface, publisher publisher.PublisherInterface) Repositorier {
	return &repository{
		db:        db,
		midtrans:  midtrans,
		publisher: publisher,
	}
}

func (repo *repository) ApprovePayment(orderID string) (res *coreapi.ChargeResponse, err error) {
	return repo.midtrans.ApprovePayment(orderID)
}

func (repo *repository) CreatePaymentLog(req model.PaymentLogRequest) (res *snap.Response, err error) {
	// prepare midtrans request data
	snapReq := &snap.Request{
		CreditCard: &snap.CreditCardDetails{
			Secure: false,
		},
		EnabledPayments: []snap.SnapPaymentType{
			snap.PaymentTypeAlfamart,
			snap.PaymentTypeAkulaku,
			snap.PaymentTypeBCAKlikpay,
			snap.PaymentTypeBRIEpay,
			snap.PaymentTypeGopay,
			snap.PaymentTypeIndomaret,
			snap.PaymentTypeMandiriEcash,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  fmt.Sprint(req.OrderID),
			GrossAmt: req.TotalPayment,
		},
	}

	// send data to midtrans
	res, err = repo.midtrans.CreateTransaction(snapReq)
	if err != nil {
		return res, fmt.Errorf("error midtrans : %v", err.Error())
	}

	// publish data to RabbitMQ
	err = repo.publisher.Publish(req, "create_payment_logs")
	if err != nil {
		err = fmt.Errorf("error publish data to RabbitMQ : %s", err.Error())
		return
	}

	return
}