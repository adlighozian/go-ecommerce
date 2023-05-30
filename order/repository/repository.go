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

func (repo *repository) CreateOrders(req []model.OrderReq) ([]model.Orders, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	var sent []model.Orders

	for _, v := range req {
		inRandom := model.Orders{
			UserID:      v.UserID,
			ShippingID:  v.ShippingID,
			TotalPrice:  v.TotalPrice,
			Status:      false,
			OrderNumber: random.NewRandom().RandomString(),
		}

		sent = append(sent, inRandom)
	}

	err := repo.sent.Public(sent, "create_order")
	if err != nil {
		return nil, errors.New("failed publisher")
	}

	time.Sleep(1 * time.Second)

	var resultss []model.Orders
	query := `select * from orders where order_number = $1`
	for _, v := range sent {

		stmt, err := repo.db.PrepareContext(ctx, query)
		failerror.FailError(err, "error prepare")

		result, err := stmt.QueryContext(ctx, v.OrderNumber)
		failerror.FailError(err, "error prepare")

		var temp model.Orders
		for result.Next() {
			result.Scan(&temp.Id, &temp.UserID, &temp.ShippingID, &temp.TotalPrice, &temp.Status, &temp.CreatedAt, &temp.UpdatedAt, &temp.OrderNumber)
		}
		if temp.Id == 0 {
			continue
		}
		resultss = append(resultss, temp)
	}

	if resultss == nil {
		return nil, errors.New("error create product")
	}

	return resultss, nil
}

func (repo *repository) ShowOrders(req model.OrderItems) (model.Orders, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	queryDetail := `select * from orders where user_id = $1 and order_number = $2`

	rows, err := repo.db.QueryContext(ctx, queryDetail, req.UserId, req.OrderNumber)
	failerror.FailError(err, "error query")

	var temp model.Orders
	for rows.Next() {
		err := rows.Scan(&temp.Id, &temp.UserID, &temp.ShippingID, &temp.TotalPrice, &temp.Status, &temp.CreatedAt, &temp.UpdatedAt, &temp.OrderNumber)
		failerror.FailError(err, "error scan")

	}

	if temp.Id == 0 {
		return model.Orders{}, errors.New("order not found")
	}

	return temp, nil
}

func (repo *repository) UpdateOrders(idOrder int, req string) error {
	return nil
}
