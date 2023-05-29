package repository

import (
	"context"
	"database/sql"
	"time"
	"review-go/model"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repositorier {
	return &repository{
		db: db,
	}
}

func (repo *repository) GetByProductID(productID int) (res []model.Review, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, user_id, product_id, rating, review_text quantity FROM reviews WHERE product_id = ?`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	result, err := stmt.QueryContext(ctx, productID)
	if err != nil {
		return
	}

	for result.Next() {
		var temp model.Review
		result.Scan(&temp.Id, &temp.UserID, &temp.ProductID, &temp.Rating, &temp.ReviewText)
		res = append(res, temp)
	}

	return
}

func (repo *repository) Create(req []model.ReviewRequest) (res []model.Review, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := `INSERT INTO reviews (user_id, product_id, rating, review_text) value (?, ?, ?, ?)`
	trx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	stmt, err := trx.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	for _, v := range req {
		result, err := stmt.ExecContext(ctx, v.UserID, v.ProductID, v.Rating, v.ReviewText)
		if err != nil {
			trx.Rollback()
			return []model.Review{}, err
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return []model.Review{}, err
		}

		res = append(res, model.Review{
			Id:   		int(lastID),
			UserID: 	v.UserID,
			ProductID: 	v.ProductID,
			Rating: 	v.Rating,
			ReviewText: v.ReviewText,
		})
	}

	trx.Commit()

	return
}