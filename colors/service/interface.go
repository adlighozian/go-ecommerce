package service

import "product-colors-go/model"

type Servicer interface {
	GetColors() (model.Respon, error)
	CreateColors(req []model.ColorsReq) (model.Respon, error)
	DeleteColors(idColors int) (model.Respon, error)
}
