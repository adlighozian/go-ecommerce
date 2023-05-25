package repositoy

import (
	"cart-go/model"
)

type Repositorier interface {
	Get() (res []model.Cart, err error)
	GetDetail(id int) (res model.Cart, err error)
	Create(req []model.CartRequest) (err error)
	Delete(id int) (err error)
}