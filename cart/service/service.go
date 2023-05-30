package service

import (
	"cart-go/model"
	"cart-go/repository"
	"errors"
	"fmt"
)

type service struct {
	repo repository.Repositorier
}

func NewService(repo repository.Repositorier) Servicer {
	return &service{
		repo: repo,
	}
}

func (svc *service) Get(userID int) (res []model.Cart, err error) {
	return svc.repo.Get(userID)
}

func (svc *service) GetDetail(userID, productID int) (res model.Cart, err error) {
	res, err = svc.repo.GetDetail(userID, productID)
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
	// check in each userID, ProductID already exist in cart or not
	for _, v := range req {
		_, err := svc.GetDetail(v.UserID, v.ProductID)
		if err != nil {
			continue
		} else {
			err = fmt.Errorf("cart with product_id %d in user_id %d already exist", v.ProductID, v.UserID)
			return []model.Cart{}, err
		}
	}

	return svc.repo.Create(req)
}

func (svc *service) Delete(cartID int) (err error) {
	// check cart id exist or not
	emptyStruct := model.Cart{}
	res, _ := svc.repo.GetByID(cartID)
	if res == emptyStruct {
		return fmt.Errorf("item with id %d not found", cartID)
	}
	return svc.repo.Delete(cartID)
}