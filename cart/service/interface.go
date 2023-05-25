package service

import (
	"cart-go/model"
)

type Servicer interface {
	Get() (res []model.Cart, err error)
	GetDetail(id int) (res model.Cart, err error)
	Create(req []model.CartRequest) (res []model.Cart, err error)
	Delete(id int) (err error)
}