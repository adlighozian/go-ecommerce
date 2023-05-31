package repository

import (
	"database/sql"
	"errors"
	"order-go/helper/failerror"
	"order-go/helper/random"
	"order-go/helper/timeout"
	"order-go/model"
	"order-go/publisher"
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

func (repo *repository) GetOrders(idUser int) ([]model.Orders, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	querySelect := `select * from orders where user_id = $1`

	stmt, err := repo.db.PrepareContext(ctx, querySelect)
	failerror.FailError(err, "error prepare context")

	result, err := stmt.QueryContext(ctx, idUser)
	if err != nil {
		return []model.Orders{}, errors.New("error get data")
	}

	var data []model.Orders

	for result.Next() {
		var temp model.Orders
		result.Scan(&temp.Id, &temp.UserID, &temp.ShippingID, &temp.TotalPrice, &temp.Status, &temp.CreatedAt, &temp.UpdatedAt, &temp.OrderNumber)
		data = append(data, temp)
	}

	return data, nil
}

func (repo *repository) CreateOrders(req model.GetOrders) (model.GetOrdersSent, error) {
	// ctx, cancel := timeout.NewCtxTimeout()
	// defer cancel()


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
		return model.GetOrdersSent{}, errors.New("failed publisher")
	}

	return inRandom, nil

	// time.Sleep(1 * time.Second)

	// var resultss []model.Orders
	// query := `select * from orders where order_number = $1`
	// for _, v := range sent {

	// 	stmt, err := repo.db.PrepareContext(ctx, query)
	// 	failerror.FailError(err, "error prepare")

	// 	result, err := stmt.QueryContext(ctx, v.OrderNumber)
	// 	failerror.FailError(err, "error prepare")

	// 	var temp model.Orders
	// 	for result.Next() {
	// 		result.Scan(&temp.Id, &temp.UserID, &temp.ShippingID, &temp.TotalPrice, &temp.Status, &temp.CreatedAt, &temp.UpdatedAt, &temp.OrderNumber)
	// 	}
	// 	if temp.Id == 0 {
	// 		continue
	// 	}
	// 	resultss = append(resultss, temp)
	// }

	// if resultss == nil {
	// 	return nil, errors.New("error create product")
	// }

	// return resultss, nil
}

func (repo *repository) ShowOrders(req model.OrderItems) ([]model.OrderItem, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	queryDetail := `select oi.id, oi.order_id, oi.product_id, oi.quantity, oi.total_price, oi.created_at, oi.updated_at from order_items oi join orders o on o.id = oi.order_id where o.order_number = $1 and o.user_id = $2`
	// queryDetail := `select * from orders where user_id = $1 and order_number = $2`

	rows, err := repo.db.QueryContext(ctx, queryDetail, req.OrderNumber, req.UserId)
	failerror.FailError(err, "error query")

	var data []model.OrderItem
	for rows.Next() {
		var temp model.OrderItem
		err := rows.Scan(&temp.Id, &temp.OrderID, &temp.ProductId, &temp.Quantity, &temp.TotalPrice, &temp.CreatedAt, &temp.UpdatedAt)
		failerror.FailError(err, "error scan")

		data = append(data, temp)
	}

	if data == nil {
		return []model.OrderItem{}, errors.New("order not found")
	}

	return data, nil
}

func (repo *repository) UpdateOrders(idOrder int, req string) error {
	return nil
}
