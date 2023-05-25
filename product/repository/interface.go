package repository

import "product-go/model"

type Repositorier interface {
	GetProduct(req model.ProductSearch) ([]model.Product, error)
	ShowProduct(id int) (model.Product, error)
	CreateProduct(req model.ProductReq) ([]model.Product, error)
	UpdateProduct(id int) error
	DeleteProduct(id int) error
}
