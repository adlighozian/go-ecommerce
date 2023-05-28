package repository

import (
	"cart-go/model"
	"cart-go/publisher"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
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

func (repo *repository) Get(userID int) (res []model.Cart, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, user_id, product_id, quantity FROM Carts WHERE user_id = $1`
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
		result.Scan(&temp.Id, &temp.UserID, &temp.ProductID, &temp.Quantity)
		res = append(res, temp)
	}

	return
}

func (repo *repository) GetDetail(userID, productID int) (res model.Cart, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, user_id, product_id, quantity FROM carts WHERE user_id = $1 AND product_id = $2`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	result, err := stmt.QueryContext(ctx, userID, productID)
	if err != nil {
		return
	}

	for result.Next() {
		result.Scan(&res.Id, &res.UserID, &res.ProductID, &res.Quantity)
	}
	return
}

func (repo *repository) Create(req []model.CartRequest) (res []model.Cart, err error) {
	// publish data to RabbitMQ
	err = repo.publisher.Publish(req, "create_carts")
	if err != nil {
		err = fmt.Errorf("error publish data to RabbitMQ : %s", err.Error())
		return
	}

	time.Sleep(3*time.Second)

	for _, v := range req {
		result, err := repo.GetDetail(v.UserID, v.ProductID)
		if err != nil {
			return []model.Cart{}, errors.New("error get by user id after create")
		}
		res = append(res, result)
	}
	return
}

func (repo *repository) Delete(cartID int) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM carts WHERE id = $1`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.QueryContext(ctx, cartID)
	if err != nil {
		return
	}

	return
}
