package service

import (
	"product-go/model"
	"product-go/repository"
)

type service struct {
	repo repository.Repositorier
}

func NewService(repo repository.Repositorier) Servicer {
	return &service{
		repo: repo,
	}
}

func (svc *service) GetProduct(req model.ProductSearch) (model.Respon, error) {
	return model.Respon{}, nil
}

func (svc *service) ShowDetail(id int) (model.Respon, error) {
	return model.Respon{}, nil
}

func (svc *service) CreateProduct(req model.ProductReq) (model.Respon, error) {
	return model.Respon{}, nil
}

func (svc *service) UpdateProduct(id int) (model.Respon, error) {
	return model.Respon{}, nil
}

func (svc *service) DeleteProduct(id int) (model.Respon, error) {
	return model.Respon{}, nil
}
