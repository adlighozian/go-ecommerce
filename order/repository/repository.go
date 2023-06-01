package repository

import (
	"database/sql"
	"errors"
	"order-go/helper/failerror"
	"order-go/helper/random"
	"order-go/helper/timeout"
	"order-go/model"
	"order-go/publisher"
	"time"
)

type repository struct {
	db   *sql.DB
	sent publisher.Publisher
}

func NewRepository(db *sql.DB, sent publisher.Publisher) Repositorier {
	return &repository{
		db:   db,
		sent: sent,
	}
}

func (repo *repository) GetOrders(userID int) ([]model.Orders, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	querySelect := `select * from orders where user_id = $1`

	result, err := repo.db.QueryContext(ctx, querySelect, userID)
	if err != nil {
		return []model.Orders{}, errors.New("error get data")
	}

	var data = []model.Orders{}

	for result.Next() {
		var temp model.Orders
		result.Scan(&temp.Id, &temp.UserID, &temp.ShippingID, &temp.TotalPrice, &temp.Status, &temp.CreatedAt, &temp.UpdatedAt, &temp.OrderNumber)
		data = append(data, temp)
	}

	return data, nil
}

func (repo *repository) ShowOrders(req model.OrderItems) (model.ResultOrders, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	var result model.ResultOrders

	queryProduct := `select oi.product_id, oi.quantity, oi.total_price from order_items oi join orders o on o.id = oi.order_id where o.order_number = $1 and o.user_id = $2`

	rows, err := repo.db.QueryContext(ctx, queryProduct, req.OrderNumber, req.UserId)
	failerror.FailError(err, "error query")

	var dataProduct = []model.OrderItemReq{}
	for rows.Next() {
		var temp model.OrderItemReq
		err := rows.Scan(&temp.ProductId, &temp.Quantity, &temp.TotalPrice)
		failerror.FailError(err, "error scan")

		dataProduct = append(dataProduct, temp)
	}

	queryOrder := `select id, user_id, shipping_id , total_price ,status , order_number,created_at ,updated_at from orders o where o.order_number = $1 and o.user_id = $2`

	var dataOrder model.Orders
	repo.db.QueryRowContext(ctx, queryOrder, req.OrderNumber, req.UserId).Scan(&dataOrder.Id, &dataOrder.UserID, &dataOrder.ShippingID, &dataOrder.TotalPrice, &dataOrder.Status, &dataOrder.OrderNumber, &dataOrder.CreatedAt, &dataOrder.UpdatedAt)

	result = model.ResultOrders{
		Id:           dataOrder.Id,
		UserID:       dataOrder.UserID,
		ShippingID:   dataOrder.ShippingID,
		TotalPrice:   dataOrder.TotalPrice,
		Status:       dataOrder.Status,
		OrderNumber:  dataOrder.OrderNumber,
		CreatedAt:    dataOrder.CreatedAt,
		UpdatedAt:    dataOrder.UpdatedAt,
		OrderItemReq: dataProduct,
	}

	if result.Id == 0 {
		return result, errors.New("order not found")
	}

	return result, nil
}

func (repo *repository) CreateOrders(req model.GetOrders) (model.Orders, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	inRandom := model.GetOrdersSent{
		UserID:       req.UserID,
		ShippingID:   req.ShippingID,
		TotalPrice:   req.TotalPrice,
		Status:       false,
		OrderNumber:  random.NewRandom().RandomString(),
		OrderItemReq: req.OrderItemReq,
	}

	err := repo.sent.Public(inRandom, "create_order")
	if err != nil {
		return model.Orders{}, errors.New("failed publisher")
	}

	time.Sleep(3 * time.Second)

	query := `select id,user_id ,shipping_id ,total_price ,status ,order_number ,created_at ,updated_at  from orders where order_number = $1`

	var temp model.Orders
	repo.db.QueryRowContext(ctx, query, inRandom.OrderNumber).Scan(&temp.Id, &temp.UserID, &temp.ShippingID, &temp.TotalPrice, &temp.Status, &temp.OrderNumber, &temp.CreatedAt, &temp.UpdatedAt)

	if temp.Id == 0 {
		return temp, errors.New("error create order")
	}

	return temp, nil
}

func (repo *repository) UpdateOrders(req model.OrderUpd) (model.Orders, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	err := repo.sent.Public(req, "update_order")
	if err != nil {
		return model.Orders{}, errors.New("failed publisher")
	}

	time.Sleep(3 * time.Second)

	queryOrder := `select id, user_id, shipping_id , total_price ,status , order_number,created_at ,updated_at from orders o where o.order_number = $1`

	var temp model.Orders
	repo.db.QueryRowContext(ctx, queryOrder, req.OrderNumber).Scan(&temp.Id, &temp.UserID, &temp.ShippingID, &temp.TotalPrice, &temp.Status, &temp.OrderNumber, &temp.CreatedAt, &temp.UpdatedAt)

	if temp.Id == 0 {
		return model.Orders{}, errors.New("error update order")
	}

	return temp, nil
}
