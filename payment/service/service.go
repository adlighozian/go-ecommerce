package service

import (
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"payment-go/model"
	"payment-go/repository"
)

type service struct {
	repo repository.Repositorier
}

func NewService(repo repository.Repositorier) Servicer {
	return &service{
		repo: repo,
	}
}

func (svc *service) CheckPayment(orderID string) (res *coreapi.TransactionStatusResponse, err error) {
	return svc.repo.CheckPayment(orderID)
}

func (svc *service) CreatePaymentLog(req model.PaymentLogRequest) (res *snap.Response, err error) {
	return svc.repo.CreatePaymentLog(req)
}