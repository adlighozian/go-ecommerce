package repository

import (
	"cart-go/model"
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

func (repo *repository) Get(userID int) (res []model.Cart, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, user_id, product_id, quantity FROM Carts WHERE user_id = ?`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	result, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return
	}

	for result.Next() {
		var temp model.Cart
		result.Scan(&temp.Id, &temp.UserID, &temp.ProductID)
		res = append(res, temp)
	}

	return
}

func (repo *repository) GetDetail(userID, cartID int) (res model.Cart, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, user_id, product_id, quantity FROM Carts WHERE user_id = ? AND WHERE cart_id = ?`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	result, err := stmt.QueryContext(ctx, userID, cartID)
	if err != nil {
		return
	}

	for result.Next() {
		var res model.Cart
		result.Scan(&res.Id, &res.UserID, &res.ProductID)
	}

	return
}

func (repo *repository) Create(req []model.CartRequest) (res []model.Cart, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := `INSERT INTO Carts (user_id, product_id, quantity) value (?, ?, ?)`
	trx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	stmt, err := trx.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	for _, v := range req {
		result, err := stmt.ExecContext(ctx, v.UserID, v.ProductID, v.Quantity)
		if err != nil {
			trx.Rollback()
			return []model.Cart{}, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return []model.Cart{}, err
		}

		res = append(res, model.Cart{
			Id:   		int(lastID),
			UserID: 	v.UserID,
			ProductID: 	v.ProductID,
			Quantity:	v.Quantity,
		})
	}

	trx.Commit()

	return
}

func (repo *repository) Delete(userID, cartID int) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM Carts WHERE user_id = ? AND WHERE cart_id = ?`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.ExecContext(ctx, userID, cartID)
	if err != nil {
		return
	}

	return
}