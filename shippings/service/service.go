package service

import (
	"errors"
	"net/http"
	"shippings-go/model"
	"shippings-go/repository"
)

type service struct {
	repo repository.Repositorier
}

func NewService(repo repository.Repositorier) Servicer {
	return &service{
		repo: repo,
	}
}

func (svc *service) GetShipping() (model.Respon, error) {

	// start
	res, err := svc.repo.GetShipping()
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

func (svc *service) CreateShipping(req []model.ShippingReq) (model.Respon, error) {

	var check []model.ShippingReq

	for _, v := range req {
		if v.Name == "" {
			continue
		}
		data := model.ShippingReq{
			Name: v.Name,
		}
		check = append(check, data)
	}

	if check == nil {
		return model.Respon{
			Status: http.StatusBadRequest,
			Data:   nil,
		}, errors.New("error input")
	}

	// start
	res, err := svc.repo.CreateShipping(check)
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

func (svc *service) DeleteShipping(id int) (model.Respon, error) {
	if id <= 0 {
		return model.Respon{
			Status: http.StatusBadRequest,
			Data:   nil,
		}, errors.New("error invalid id")
	}

	// start
	res, err := svc.repo.DeleteShipping(id)
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
