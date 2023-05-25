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

func (repo *repository) GetDetail(id int) (res []model.Review, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, user_id, product_id, quantity FROM reviews WHERE = ?`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err = stmt.QueryContext(ctx, id)
	if err != nil {
		return
	}

	for res.Next() {
		var temp model.Review
		res.Scan(&temp.Id, &temp.UserID, &temp.ProductID, &temp.Quantity)
		res = append(res, temp)
	}

	return
}

func (repo *repository) Create(req []model.ReviewRequest) (res []model.Review, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := `INSERT INTO reviews (user_id, product_id, quantity) value (?, ?,?)`
	trx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	for _, v := range req {
		result, err := stmt.ExecContext(ctx, v.UserID, v.ProductID, v.Quantity)
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
			Quantity: 	v.Quantity,
		})
	}

	trx.Commit()

	return
}