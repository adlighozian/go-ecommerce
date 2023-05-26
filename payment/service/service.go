package service

import (
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

func (svc *service) GetPaymentMethod() (res []model.PaymentMethod, err error) {
	return svc.repo.GetPaymentMethod()
}

func (svc *service) CreatePaymentMethod(req []model.PaymentMethodRequest) (res []model.PaymentMethod, err error) {
	return svc.repo.CreatePaymentMethod(req)
}

func (svc *service) CreatePaymentLog(req []model.PaymentLogsRequest) (res []model.PaymentLogsRequest, err error) {
	return svc.repo.CreatePaymentLog(req)
}