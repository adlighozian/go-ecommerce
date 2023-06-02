package repository

import (
	"consumer-insert-voucher-go/helpers"
	"consumer-insert-voucher-go/model"
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

func (p product) CreateProduct(req []model.Voucher) error {
	ctx, cancel := helpers.NewCtxTimeout()
	defer cancel()

	trx, err := p.db.BeginTx(ctx, nil)
	helpers.FailOnError(err, "error config")

	query := `insert into voucher (store_id,product_id,category_id,discount_value,name,code,start_date,end_date) 
	values ($1,	$2,$3,$4,$5,$6,$7,$8) returning id`

	stmt, err := trx.PrepareContext(ctx, query)
	if err != nil {
		trx.Rollback()
		return err
	}

	var ids []int
	for _, v := range req {
		var id int

		err = stmt.QueryRowContext(ctx, v.StoreID, v.ProductID, v.CategoryID, v.Discount, v.Name, v.Code, v.StartDate, v.EndDate).Scan(&id)
		if err != nil {
			trx.Rollback()
		}
		ids = append(ids, id)
	}
	trx.Commit()

	fmt.Println(ids)

	return nil

}
