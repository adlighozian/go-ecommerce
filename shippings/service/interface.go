package service

import "shippings-go/model"

type Servicer interface {
	GetShipping() (model.Respon, error)
	CreateShipping(req []model.ShippingReq) (model.Respon, error)
	DeleteShipping(idShipping int) (model.Respon, error)
}
