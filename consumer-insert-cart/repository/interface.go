package repository

import (
	"consumer-cart-go/model"
)

type Repositorier interface {
	Create(req []model.CartRequest) (err error)
}