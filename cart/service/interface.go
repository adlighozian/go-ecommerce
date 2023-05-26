package service

import (
	"cart-go/model"
)

type Servicer interface {
	Get(userID int) (res []model.Cart, err error)
	GetDetail(userID, cartID int) (res model.Cart, err error)
	Create(req []model.CartRequest) (res []model.Cart, err error)
	Delete(userID, cartID int) (err error)
}