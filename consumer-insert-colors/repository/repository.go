package repository

import (
	"consumer-insert-shipping-go/helpers"
	"consumer-insert-shipping-go/model"
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

func (p product) CreateProduct(req []model.ColorsReq) error {
	ctx, cancel := helpers.NewCtxTimeout()
	defer cancel()

	trx, err := p.db.BeginTx(ctx, nil)
	helpers.FailOnError(err, "error config")

	query := `insert into product_colors (name) values ($1) returning id`

	stmt, err := trx.PrepareContext(ctx, query)
	if err != nil {
		trx.Rollback()
		return err
	}

	var ids []int
	for _, v := range req {
		var id int
		err = stmt.QueryRowContext(ctx, v.Name).Scan(&id)
		if err != nil {
			trx.Rollback()
		}
		ids = append(ids, id)
	}
	trx.Commit()

	fmt.Println(ids)

	return nil

}
