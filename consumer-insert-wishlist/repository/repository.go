package repository

import (
	"consumer-wishlist-go/model"
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

func (repo *repository) Create(req []model.WishlistRequest) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := `INSERT INTO wishlists (user_id, product_id) values ($1, $2)`
	trx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	stmt, err := trx.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	for _, v := range req {
		_, err := stmt.ExecContext(ctx, v.UserID, v.ProductID)
		if err != nil {
			trx.Rollback()
			return err
		}
	}
	
	trx.Commit()
	return
}