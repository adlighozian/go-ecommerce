package repository

import (
	"consumer-insert-order-go/helpers"
	"consumer-insert-order-go/model"
	"database/sql"
	"fmt"
)

type Product interface {
	CreateProduct(req []model.Orders) error
}

type product struct {
	db *sql.DB
}

func NewProduct(db *sql.DB) Product {
	return &product{
		db: db,
	}
}

func (p product) CreateProduct(req []model.Orders) error {
	ctx, cancel := helpers.NewCtxTimeout()
	defer cancel()

	trx, err := p.db.BeginTx(ctx, nil)
	helpers.FailOnError(err, "error config")

	query := `insert into orders (user_id, shipping_id, total_price, status, order_number) values ($1,$2,$3,$4,$5) returning id`

	var ids []int
	for _, v := range req {
		var id int

		stmt, err := trx.PrepareContext(ctx, query)
		if err != nil {
			trx.Rollback()
			return err
		}

		err = stmt.QueryRowContext(ctx, v.UserID, v.ShippingID, v.TotalPrice, v.Status, v.OrderNumber).Scan(&id)
		if err != nil {
			trx.Rollback()
		}
		ids = append(ids, id)
	}
	trx.Commit()

	fmt.Println(ids)

	return nil

}
