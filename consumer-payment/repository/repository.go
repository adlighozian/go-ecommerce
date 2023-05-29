package repository

import (
	"consumer-payment-logs-go/model"
	"context"
	"database/sql"
	"time"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repositorier {
	return &repository{
		db: db,
	}
}

func (repo *repository) Create(req model.PaymentLogRequest) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := `INSERT INTO payment_logs (user_id, order_id, payment_method_id, total_payment) values ($1, $2, $3, $4)`
	trx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	stmt, err := trx.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.ExecContext(ctx, req.UserID, req.OrderID, req.PaymentMethodID, req.TotalPayment)
	if err != nil {
		trx.Rollback()
		return err
	}

	trx.Commit()
	return
}
