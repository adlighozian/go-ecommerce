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

func (repo *repository) GetPaymentMethod() (res []model.PaymentMethod, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, payment_gateway_id, name FROM payment_methods`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	result, err := stmt.QueryContext(ctx)
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

func (repo *repository) GetPaymentMethodByID(paymentMethodID int) (res model.PaymentMethod, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, payment_gateway_id, name FROM payment_methods WHERE id = $1`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	result, err := stmt.QueryContext(ctx, paymentMethodID)
	if err != nil {
		return
	}

	for result.Next() {
		result.Scan(&res.Id, &res.PaymentGatewayID, &res.Name)
	}

	return
}

func (repo *repository) GetPaymentMethodByName(name string) (res model.PaymentMethod, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, payment_gateway_id, name FROM payment_methods WHERE name = $1`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	result, err := stmt.QueryContext(ctx, name)
	if err != nil {
		return
	}

	for result.Next() {
		result.Scan(&res.Id, &res.PaymentGatewayID, &res.Name)
	}

	return
}

func (repo *repository) CreatePaymentMethod(req []model.PaymentMethodRequest) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := `INSERT INTO payment_methods(payment_gateway_id, name) values ($1, $2)`
	trx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	stmt, err := trx.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	for _, v := range req {
		_, err = stmt.ExecContext(ctx, v.PaymentGatewayID, v.Name)
		if err != nil {
			trx.Rollback()
			return err
		}
	}

	trx.Commit()

	return
}

func (repo *repository) CreatePaymentLog(req model.PaymentLogRequest) (res model.PaymentLog, err error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	// defer cancel()

	// query := `INSERT INTO payment_logs (user_id, order_id, payment_method_id) values ($1, $2, $3)`
	// trx, err := repo.db.BeginTx(ctx, nil)
	// if err != nil {
	// 	return
	// }

	// stmt, err := trx.PrepareContext(ctx, query)
	// if err != nil {
	// 	return
	// }

	// _, err = stmt.ExecContext(ctx, v.UserID, v.OrderID, v.PaymentMethodID)
	// if err != nil {
	// 	trx.Rollback()
	// 	return model.PaymentLog{}, err
	// }

	// 	// res = append(res, model.PaymentMethod{
	// 	// 	Id:   				int(lastID),
	// 	// 	PaymentGatewayID: 	v.PaymentGatewayID,
	// 	// 	Name: 				v.Name,
	// 	// })

	// trx.Commit()
	return
}

func (repo *repository) DeletePaymentMethod(paymentMethodID int) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM payment_methods WHERE id = $1`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.QueryContext(ctx, paymentMethodID)
	if err != nil {
		return
	}

	return
}
