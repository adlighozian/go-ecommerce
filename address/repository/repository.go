package repository

import (
	"address-go/model"
	"address-go/publisher"
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

func (repo *repository) Get(userID int) (res []model.Address, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, user_id, street, state, city, country, zipcode, phone_number FROM addresses WHERE user_id = $1`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	result, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return
	}

	for result.Next() {
		var temp model.Address
		result.Scan(&temp.Id, &temp.UserID, &temp.Street, &temp.State, &temp.City, &temp.Country, &temp.Zipcode, &temp.PhoneNumber)
		res = append(res, temp)
	}
	return
}

func (repo *repository) Create(req model.AddressRequest) (res []model.Address, err error) {
	// publish data to RabbitMQ
	err = repo.publisher.Publish(req, "create_address")
	if err != nil {
		err = fmt.Errorf("error publish data to RabbitMQ : %s", err.Error())
		return
	}

	time.Sleep(3 * time.Second)

	res, err = repo.Get(req.UserID)
	if err != nil {
		return []model.Address{}, errors.New("error get addresses by user id after create")
	}
	return
}

func (repo *repository) Delete(addressID int) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM addresses WHERE id = $1`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.QueryContext(ctx, addressID)
	if err != nil {
		return
	}

	return
}
