package repository

import (
	"consumer-cart-go/model"
	"context"
	"database/sql"
	"time"
)

type repository struct {
	db 		  *sql.DB
}

func NewRepository(db *sql.DB) Repositorier {
	return &repository{
		db: db,
	}
}

func (repo *repository) Create(req []model.CartRequest) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := `INSERT INTO carts (user_id, product_id, quantity) values ($1, $2, $3) returning id, user_id, product_id, quantity`
	trx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	stmt, err := trx.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	for _, v := range req {
		_, err := stmt.ExecContext(ctx, v.UserID, v.ProductID, v.Quantity)
		if err != nil {
			trx.Rollback()
			return err
		}
	}
	
	trx.Commit()
	return
}