package repository

import (
	"consumer-product-go/helpers"
	"consumer-product-go/model"
	"database/sql"
	"fmt"
)

type product struct {
	db *sql.DB
}

func NewProduct(db *sql.DB) Product {
	return &product{
		db: db,
	}
}

func (p product) UpdateProduct(req model.ProductUpd) error {
	ctx, cancel := helpers.NewCtxTimeout()
	defer cancel()

	querys := `update products set store_id = $1 ,category_id = $2, size_id = $3, color_id = $4, name = $5, subtitle = $6,description = $7, unit_price = $8, status = $9, stock = $10, weight = $11, brand = $12 where id = $13`

	_, err := p.db.ExecContext(ctx, querys, req.StoreID, req.CategoryID, req.SizeID, req.ColorID, req.Name, req.Subtitle, req.Description, req.UnitPrice, req.Status, req.Stock, req.Weight, req.Brand, req.Id)
	helpers.FailOnError(err, "error exec")

	fmt.Println(req.Id)

	return nil

}
