package service

import (
	"errors"
	"fmt"
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

func (svc *service) GetPaymentMethodByName(name string) (res model.PaymentMethod, err error) {
	res, err = svc.repo.GetPaymentMethodByName(name)
	if err != nil {
		return
	}

	emptyStruct := model.PaymentMethod{}
	if res == emptyStruct {
		err = errors.New("payment method not found")
		return
	}
	return
}

func (svc *service) CreatePaymentMethod(req []model.PaymentMethodRequest) (res []model.PaymentMethod, err error) {
	// check if payment method already exist or not
	for _, v := range req {
		emptyStruct := model.PaymentMethod{}
		res, _ := svc.GetPaymentMethodByName(v.Name)
		if res != emptyStruct {
			return []model.PaymentMethod{}, fmt.Errorf("payment method %s already exist", v.Name)
		} else {
			continue
		}
	}
	
	err = svc.repo.CreatePaymentMethod(req)
	if err != nil {
		return
	}
	// get data after insert, to check if data already inserted in database
	for _, v := range req {
		paymentMethod, err := svc.GetPaymentMethodByName(v.Name)
		if err != nil {
			return []model.PaymentMethod{}, fmt.Errorf("error get data after create: %s", err.Error())
		}

		res = append(res, model.PaymentMethod{
			Id:   				paymentMethod.Id,
			PaymentGatewayID: 	v.PaymentGatewayID,
			Name: 				v.Name,
		})
	}
	return
}

func (svc *service) CreatePaymentLog(req model.PaymentLogRequest) (res model.PaymentLog, err error) {
	return svc.repo.CreatePaymentLog(req)
}

func (svc *service) DeletePaymentMethod(paymentMethodID int) (err error) {
	// check cart id exist or not
	emptyStruct := model.PaymentMethod{}
	res, _ := svc.repo.GetPaymentMethodByID(paymentMethodID)
	if res == emptyStruct {
		return fmt.Errorf("payment method with id %d not found", paymentMethodID)
	}
	return svc.repo.DeletePaymentMethod(paymentMethodID)
}