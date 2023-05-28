package repository

import (
	"cart-go/model"
)

type Repositorier interface {
	Get(userID int) (res []model.Cart, err error)
	GetDetail(userID, productID int) (res model.Cart, err error)
	Create(req []model.CartRequest) (res []model.Cart, err error)
	Delete(cartID int) (err error)
}