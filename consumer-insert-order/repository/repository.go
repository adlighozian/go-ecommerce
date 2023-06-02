package repository

import (
	"consumer-insert-order-go/helpers"
	"consumer-insert-order-go/model"
	"database/sql"
	"fmt"
)

type Product interface {
	CreateProduct(req model.GetOrdersSent) error
}

type product struct {
	db *sql.DB
}

func NewProduct(db *sql.DB) Product {
	return &product{
		db: db,
	}
}

func (p product) CreateProduct(req model.GetOrdersSent) error {
	ctx, cancel := helpers.NewCtxTimeout()
	defer cancel()

	trx, err := p.db.BeginTx(ctx, nil)
	helpers.FailOnError(err, "error config")

	queryOrder := `insert into orders (user_id, shipping_id, total_price, status, order_number) values ($1,$2,$3,$4,$5) returning id`

	var idOrders int
	err = trx.QueryRowContext(ctx, queryOrder, req.UserID, req.ShippingID, req.TotalPrice, req.Status, req.OrderNumber).Scan(&idOrders)
	if err != nil {
		trx.Rollback()
	}

	queryOrderItem := `insert into order_items (order_id,product_id,quantity,total_price) values ($1,$2,$3,$4)`

	stmt, err := trx.PrepareContext(ctx, queryOrderItem)
	if err != nil {
		trx.Rollback()
	}

	for _, v := range req.OrderItemReq {
		_, err := stmt.ExecContext(ctx, idOrders, v.ProductId, v.Quantity, v.TotalPrice)
		if err != nil {
			trx.Rollback()
		}
	}

	trx.Commit()

	fmt.Println(idOrders)

	return nil

}
