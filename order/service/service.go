package service

import (
	"errors"
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
	if idUser <= 0 {
		return model.Respon{
			Status: http.StatusBadRequest,
			Data:   nil,
		}, errors.New("invalid input id")
	}

	// start
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

func (svc *service) CreateOrders(req model.GetOrders) (model.Respon, error) {
	var check int

	if req.UserID == 0 || req.ShippingID == 0 || req.TotalPrice == 0 || req.OrderItemReq == nil || len(req.OrderItemReq) == 0 {
		check++
	}

	for _, v := range req.OrderItemReq {
		if v.ProductId <= 0 || v.Quantity <= 0 || v.TotalPrice <= 0 {
			check++
		}
	}

	if check != 0 {
		return model.Respon{
			Status: http.StatusBadRequest,
			Data:   nil,
		}, errors.New("invalid input or item null")
	}

	// start
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
	if req.OrderNumber == "" || req.UserId <= 0 {
		return model.Respon{
			Status: http.StatusBadRequest,
			Data:   nil,
		}, errors.New("invalid input")
	}

	// start
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

func (svc *service) UpdateOrders(req model.OrderUpd) (model.Respon, error) {

	if req.OrderNumber == "" {
		return model.Respon{
			Status: http.StatusBadRequest,
			Data:   nil,
		}, errors.New("invalid input")
	}

	// start
	res, err := svc.repo.UpdateOrders(req)
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
