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

func (p product) CreateProduct(req []model.Product) error {
	ctx, cancel := helpers.NewCtxTimeout()
	defer cancel()

	trx, err := p.db.BeginTx(ctx, nil)
	helpers.FailOnError(err, "error config")

	query := `insert into products (store_id,category_id,size_id,color_id,name,subtitle,description,unit_price,status,stock,sku,weight,brand) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) returning id`

	stmt, err := trx.PrepareContext(ctx, query)
	if err != nil {
		trx.Rollback()
		return err
	}

	var ids []int
	for _, v := range req {
		var id int

		err = stmt.QueryRowContext(ctx, v.StoreID, v.CategoryID, v.SizeID, v.ColorID, v.Name, v.Subtitle, v.Description, v.UnitPrice, v.Status, v.Stock, v.Sku, v.Weight, v.Brand).Scan(&id)
		if err != nil {
			trx.Rollback()
		}
		ids = append(ids, id)
	}
	trx.Commit()

	fmt.Println(ids)

	return nil

}
