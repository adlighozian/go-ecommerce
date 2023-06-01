package service

import "product-go/model"

type Servicer interface {
	GetProduct(req model.ProductSearch) (model.Respon, error)
	ShowProduct(id int) (model.Respon, error)
	CreateProduct(req []model.ProductReq) (model.Respon, error)
	UpdateProduct(req model.ProductUpd, id int) (model.Respon, error)
	DeleteProduct(id int) (model.Respon, error)
}
