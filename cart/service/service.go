package service

import (
	"errors"
	"cart-go/model"
	"cart-go/repository"
)

type service struct {
	repo repository.Repositorier
}

func NewService(repo repository.Repositorier) Servicer {
	return &service{
		repo: repo,
	}
}

func (svc *service) Get() (res []model.Cart, err error) {
	return svc.repo.Get()
}

func (svc *service) GetDetail(id int) (res model.Cart, err error) {
	res, err = svc.repo.GetDetail(id)
	if err != nil {
		return
	}

	emptyStruct := model.Cart{}
	if res == emptyStruct {
		err = errors.New("cart not found")
		return
	}
	return
}

func (svc *service) Create(req []model.CartRequest) (res []model.Cart, err error) {
	return svc.repo.Create(req)
}

func (svc *service) Delete(id int) (err error) {
	return svc.repo.Delete(id)
}