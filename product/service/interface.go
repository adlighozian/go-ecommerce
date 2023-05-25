package service

import "product-go/model"

type Servicer interface {
	GetProduct(req model.ProductSearch) (model.Respon, error)
	ShowDetail(id int) (model.Respon, error)
	Create(req model.ProductReq) (model.Respon, error)
	Update(id int) (model.Respon, error)
	Delete(id int) (model.Respon, error)
}
