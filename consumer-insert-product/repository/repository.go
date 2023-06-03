package repository

import (
	"consumer-product-go/helpers"
	"consumer-product-go/model"
	"database/sql"
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

	queryInsert := `insert into products (store_id,category_id,size_id,color_id,name,subtitle,description,unit_price,status,stock,sku,weight,brand) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) returning id`
	queryImage := `insert into image (name,image_url) values ($1,$2) returning id`
	queryMappingImage := `insert into image_products (product_id, image_id) values ($1,$2)`

	stmtInsert, err := trx.PrepareContext(ctx, queryInsert)
	helpers.FailOnError(err, "error prepare")
	stmtImage, err := trx.PrepareContext(ctx, queryImage)
	helpers.FailOnError(err, "error prepare")
	stmtMappingImage, err := trx.PrepareContext(ctx, queryMappingImage)
	helpers.FailOnError(err, "error prepare")

	// var idProduct []int
	for _, v := range req {
		var idPro int
		err = stmtInsert.QueryRowContext(ctx, v.StoreID, v.CategoryID, v.SizeID, v.ColorID, v.Name, v.Subtitle, v.Description, v.UnitPrice, v.Status, v.Stock, v.Sku, v.Weight, v.Brand).Scan(&idPro)
		if err != nil {
			trx.Rollback()
		}

		for _, v := range v.ProductImage {
			var idImag int
			err = stmtImage.QueryRowContext(ctx, v.Name, v.ImageURL).Scan(&idImag)
			if err != nil {
				trx.Rollback()
			}

			_, err = stmtMappingImage.ExecContext(ctx, idPro, idImag)
			if err != nil {
				trx.Rollback()
			}
		}

		// idProduct = append(idProduct, idPro)
	}
	trx.Commit()

	return nil

}
