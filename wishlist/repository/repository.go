package repository

import (
	"context"
	"database/sql"
	"time"
	"wishlist-go/model"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repositorier {
	return &repository{
		db: db,
	}
}

func (repo *repository) Get(userID int) (res []model.Wishlist, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, user_id, product_id FROM wishlists WHERE user_id = ?`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	result, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return
	}

	for result.Next() {
		var temp model.Wishlist
		result.Scan(&temp.Id, &temp.UserID, &temp.ProductID)
		res = append(res, temp)
	}

	return
}

func (repo *repository) GetDetail(userID, wishlistID int) (res model.Wishlist, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, user_id, product_id FROM wishlists WHERE userID = ? AND WHERE wishlist_id = ?`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	result, err := stmt.QueryContext(ctx, userID, wishlistID)
	if err != nil {
		return
	}

	for result.Next() {
		var res model.Wishlist
		result.Scan(&res.Id, &res.UserID, &res.ProductID)
	}

	return
}

func (repo *repository) Create(req []model.WishlistRequest) (res []model.Wishlist, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := `INSERT INTO wishlists (user_id, product_id) value (?, ?)`
	trx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	stmt, err := trx.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	for _, v := range req {
		result, err := stmt.ExecContext(ctx, v.UserID, v.ProductID)
		if err != nil {
			trx.Rollback()
			return []model.Wishlist{}, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return []model.Wishlist{}, err
		}

		res = append(res, model.Wishlist{
			Id:   		int(lastID),
			UserID: 	v.UserID,
			ProductID: 	v.ProductID,
		})
	}

	trx.Commit()

	return
}

func (repo *repository) Delete(userID, wishlistID int) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM wishlists WHERE id = ?`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.ExecContext(ctx, userID, wishlistID)
	if err != nil {
		return
	}

	return
}