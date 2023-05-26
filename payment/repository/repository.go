package repository

import (
	"context"
	"database/sql"
	"time"
	"payment-go/model"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repositorier {
	return &repository{
		db: db,
	}
}

func (svc *service) GetPaymentMethod() (res []model.PaymentMethod, err error) { {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, payment_gateway_id, name FROM payment_methods`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	result, err := stmt.QueryContext(ctx, productID)
	if err != nil {
		return
	}

	for result.Next() {
		var temp model.PaymentMethod
		result.Scan(&temp.Id, &temp.PaymentGatewayID, &temp.Name)
		res = append(res, temp)
	}

	return
}

func (repo *repository) CreatePaymentMethod(req []model.PaymentMethodRequest) (res []model.PaymentMethod, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := `INSERT INTO payment_methods(payment_gateway_id, name) value (?, ?)`
	trx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	stmt, err := trx.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	for _, v := range req {
		result, err := stmt.ExecContext(ctx, v.PaymentGatewayID, v.Name)
		if err != nil {
			trx.Rollback()
			return []model.Review{}, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return []model.Review{}, err
		}

		res = append(res, model.PaymentMethod{
			Id:   				int(lastID),
			PaymentGatewayID: 	v.PaymentGatewayID,
			Name: 				v.Name,
		})
	}

	trx.Commit()

	return
}

func (repo *repository) CreatePaymentLogs(req []model.PaymentLogsRequest) (res []model.PaymentLogs, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := `INSERT INTO payment_logs(user_id, order_id, payment_method_id) value (?, ?, ?)`
	trx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	stmt, err := trx.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	for _, v := range req {
		result, err := stmt.ExecContext(ctx, v.UserID, v.OrderID, v.PaymentMethod)
		if err != nil {
			trx.Rollback()
			return []model.Review{}, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return []model.Review{}, err
		}

		res = append(res, model.PaymentMethod{
			Id:   				int(lastID),
			PaymentGatewayID: 	v.PaymentGatewayID,
			Name: 				v.Name,
		})
	}

	trx.Commit()

	return
}