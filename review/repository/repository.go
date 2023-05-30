package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"review-go/model"
	"review-go/publisher"
)

type repository struct {
	db        *sql.DB
	publisher publisher.PublisherInterface
}

func NewRepository(db *sql.DB, publisher publisher.PublisherInterface) Repositorier {
	return &repository{
		db:        db,
		publisher: publisher,
	}
}

func (repo *repository) GetByProductID(productID int) (res []model.Review, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, user_id, product_id, rating, review_text quantity FROM reviews WHERE product_id = $1`
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

func (repo *repository) GetDetail(userID, productID int) (res model.Review, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, user_id, product_id, quantity FROM reviews WHERE user_id = $1 AND product_id = $2`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	result, err := stmt.QueryContext(ctx, userID, productID)
	if err != nil {
		return
	}

	for result.Next() {
		result.Scan(&res.Id, &res.UserID, &res.ProductID)
	}
	return
}

func (repo *repository) Create(req []model.ReviewRequest) (res []model.Review, err error) {
	// publish data to RabbitMQ
	err = repo.publisher.Publish(req, "create_reviews")
	if err != nil {
		err = fmt.Errorf("error publish data to RabbitMQ : %s", err.Error())
		return
	}

	time.Sleep(3*time.Second)

	for _, v := range req {
		result, err := repo.GetDetail(v.UserID, v.ProductID)
		if err != nil {
			return []model.Review{}, errors.New("error get by user id after create")
		}
		res = append(res, result)
	}
	return
}