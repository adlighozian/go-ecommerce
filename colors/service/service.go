package service

import (
	"errors"
	"net/http"
	"product-colors-go/model"
	"product-colors-go/repository"
)

type service struct {
	repo repository.Repositorier
}

func NewService(repo repository.Repositorier) Servicer {
	return &service{
		repo: repo,
	}
}

func (svc *service) GetColors() (model.Respon, error) {
	// start
	res, err := svc.repo.GetColors()
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

func (svc *service) CreateColors(req []model.ColorsReq) (model.Respon, error) {

	var check []model.ColorsReq

	for _, v := range req {
		if v.Name == "" {
			continue
		}
		data := model.ColorsReq{
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
	res, err := svc.repo.CreateColors(check)
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

func (svc *service) DeleteColors(id int) (model.Respon, error) {
	if id <= 0 {
		return model.Respon{
			Status: http.StatusBadRequest,
			Data:   nil,
		}, errors.New("error invalid id")
	}

	// start
	res, err := svc.repo.DeleteColors(id)
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
