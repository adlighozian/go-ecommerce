package repository

import (
	"product-go/model"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repositorier {
	return &repository{
		db: db,
	}
}

func (repo *repository) GetProduct(req model.ProductSearch) ([]model.Product, error) {

	// searchProduct := `select * from products p join store st on st.id = p.id join category c on c.id = p.id join sizez s on si.id = p.id join color co on co.id = p.id where st.names like '%?%' and where c.names like '%?%' and where s.names like '%?%' and where co.names like '%?%' `

	return []model.Product{}, nil
}

func (repo *repository) ShowProduct(id int) (model.Product, error) {
	return model.Product{}, nil
}

func (repo *repository) CreateProduct(req []model.ProductReq) ([]model.Product, error) {
	return []model.Product{}, nil
}

func (repo *repository) UpdateProduct(req model.ProductReq) error {
	return nil
}

func (repo *repository) DeleteProduct(id int) error {
	return nil
}
