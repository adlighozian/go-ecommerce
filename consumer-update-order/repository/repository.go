package repository

import (
	"consumer-update-order-go/helpers"
	"consumer-update-order-go/model"
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

func (p product) UpdateProduct(req model.OrderUpd) error {
	ctx, cancel := helpers.NewCtxTimeout()
	defer cancel()

	querys := `update orders set status = $1, receipt_number = $2 where order_number = $3 returning id`

	var idCheck int
	p.db.QueryRowContext(ctx, querys, req.Status, req.ReceiptNumber, req.OrderNumber).Scan(&idCheck)

	fmt.Println(idCheck)

	return nil

}
