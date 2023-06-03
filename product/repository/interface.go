package repository

import "product-go/model"

type Repositorier interface {
	GetProduct(req model.ProductSearch) ([]model.Product, error)
	ShowProduct(id int) (model.Product, error)
	CreateProduct(req []model.Product) ([]model.Product, error)
	UpdateProduct(req model.ProductUpd) (model.Product, error)
	DeleteProduct(id int) (int, error)
}
