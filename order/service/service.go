package service

import (
	"errors"
	"log"
	"net/http"
	"order-go/model"
	"order-go/repository"
)

type service struct {
	repo repository.Repositorier
}

func NewService(repo repository.Repositorier) Servicer {
	return &service{
		repo: repo,
	}
}

func (svc *service) GetOrders(idUser int) (model.Respon, error) {
	if idUser == 0 {
		return model.Respon{
			Status: http.StatusBadRequest,
			Data:   nil,
		}, errors.New("invalid input id")
	}

	res, err := svc.repo.GetOrders(idUser)
	if err != nil {
		return model.Respon{
			Status: http.StatusInternalServerError,
			Data:   nil,
		}, err
	}
	return model.Respon{
		Status: http.StatusOK,
		Data:   res,
	}, nil
}

func (svc *service) CreateOrders(req []model.OrderReq) (model.Respon, error) {

	var data []model.OrderReq

	for _, v := range req {
		if v.UserID == 0 || v.ShippingID == 0 || v.TotalPrice == 0 {
			continue
		}
		data = append(data, model.OrderReq{
			UserID:     v.UserID,
			ShippingID: v.ShippingID,
			TotalPrice: v.TotalPrice,
		})
	}

	if data == nil {
		return model.Respon{
			Status: http.StatusBadRequest,
			Data:   nil,
		}, errors.New("invalid input")
	}

	log.Println(data)

	res, err := svc.repo.CreateOrders(req)
	if err != nil {
		return model.Respon{
			Status: http.StatusInternalServerError,
			Data:   nil,
		}, err
	}
	return model.Respon{
		Status: http.StatusOK,
		Data:   res,
	}, nil
}

func (svc *service) ShowOrders(req model.OrderItems) (model.Respon, error) {

	res, err := svc.repo.ShowOrders(req)
	if err != nil {
		return model.Respon{
			Status: http.StatusInternalServerError,
			Data:   nil,
		}, err
	}
	return model.Respon{
		Status: http.StatusOK,
		Data:   res,
	}, nil
}

func (svc *service) UpdateOrders(idOrder int, req string) (model.Respon, error) {

	err := svc.repo.UpdateOrders(idOrder, req)
	if err != nil {
		return model.Respon{
			Status: http.StatusInternalServerError,
			Data:   nil,
		}, err
	}
	return model.Respon{
		Status: http.StatusOK,
		Data:   nil,
	}, nil
}
