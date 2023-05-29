package service

import (
	"order-go/model"
)

type service struct {
	repo repository.Repositorier
}

func NewService(repo repository.Repositorier) Servicer {
	return &service{
		repo: repo,
	}
}

func (svc *service) GetOrders(id int) (model.Respon, error) {
	return model.Respon{}, nil
}

func (svc *service) CreateOrders(req model.Orders) (model.Respon, error) {
	return model.Respon{}, nil
}

func (svc *service) ShowOrders(idOrder int, idUser int) (model.Respon, error) {
	return model.Respon{}, nil
}

func (svc *service) UpdateOrders(idOrder int, req string) (model.Respon, error) {
	return model.Respon{}, nil
}
