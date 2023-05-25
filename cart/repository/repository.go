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

func (repo *repository) Get() (res []model.Cart, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, user_id, product_id, quantity FROM Carts`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	result, err := stmt.QueryContext(ctx)
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

func (repo *repository) GetDetail(id int) (res model.Cart, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, user_id, product_id, quantity FROM Carts WHERE = ?`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	result, err := stmt.QueryContext(ctx, id)
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

	query := `INSERT INTO Carts (user_id, product_id, quantity) value (?, ?,?)`
	trx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	for _, v := range req {
		result, err := stmt.ExecContext(ctx, v.UserID, v.ProductID)
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
		})
	}

	trx.Commit()

	return
}

func (repo *repository) Delete(id int) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM Carts WHERE id = ?`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return
	}

	return
}