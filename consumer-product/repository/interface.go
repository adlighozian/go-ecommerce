package repository

import "consumer-product-go/model"

type Product interface {
	CreateProduct(req []model.ProductReq) error
	UpdateProduct(req model.ProductReq) error
}
