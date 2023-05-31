package repository

import (
	"database/sql"
	"fmt"
	midtransrepo "payment-go/midtrans"
	"payment-go/model"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type repository struct {
	db        *sql.DB
	midtrans  midtransrepo.MidtransInterface
}

func NewRepository(db *sql.DB, midtrans midtransrepo.MidtransInterface) Repositorier {
	return &repository{
		db:        db,
		midtrans:  midtrans,
	}
}

func (repo *repository) CheckPayment(orderID string) (res *coreapi.TransactionStatusResponse, err error) {
	return repo.midtrans.CheckPayment(orderID)
}

func (repo *repository) CreatePaymentLog(req model.PaymentLogRequest) (res *snap.Response, err error) {
	// prepare midtrans request data
	snapReq := &snap.Request{
		Metadata: model.Customer {
			UserID: req.UserID,
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
	res, err = repo.midtrans.CreatePayment(snapReq)
	if err != nil {
		return res, fmt.Errorf("error midtrans : %v", err.Error())
	}

	return
}