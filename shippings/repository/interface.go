package repository

import "shippings-go/model"

type Repositorier interface {
	GetShipping() ([]model.Shipping, error)
	CreateShipping(req []model.ShippingReq) ([]model.Shipping, error)
	DeleteShipping(id int) (int, error)
}
