package repository

import (
	"database/sql"
	"product-go/model"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repositorier {
	return &repository{
		db: db,
	}
}

func (repo *repository) GetProduct(req model.ProductSearch) ([]model.Product, error) {

	searchProduct := ""

	return []model.Product{}, nil
}

func (repo *repository) ShowProduct(id int) (model.Product, error) {
	return model.Product{}, nil
}

func (repo *repository) CreateProduct(req model.ProductReq) ([]model.Product, error) {
	return []model.Product{}, nil
}

func (repo *repository) UpdateProduct(id int) error {
	return nil
}

func (repo *repository) DeleteProduct(id int) error {
	return nil
}
