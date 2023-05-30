package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"store-go/model"
	"store-go/publisher"
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

func (repo *repository) Get() (res []model.Store, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, address_id, description, image_url, name FROM stores`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	result, err := stmt.QueryContext(ctx)
	if err != nil {
		return
	}

	for result.Next() {
		var temp model.Store
		result.Scan(&temp.Id, &temp.AddressID, &temp.Description, &temp.ImageURL, &temp.Name)
		res = append(res, temp)
	}

	return
}

func (repo *repository) GetStoreByName(name string) (res model.Store, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, address_id, description, image_url, name FROM stores WHERE name = $1`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	result, err := stmt.QueryContext(ctx, name)
	if err != nil {
		return
	}

	for result.Next() {
		result.Scan(&res.Id, &res.AddressID, &res.Description, &res.ImageURL, &res.Name)
	}
	return
}

func (repo *repository) Create(req []model.StoreRequest) (res []model.Store, err error) {
	// publish data to RabbitMQ
	err = repo.publisher.Publish(req, "create_stores")
	if err != nil {
		err = fmt.Errorf("error publish data to RabbitMQ : %s", err.Error())
		return
	}

	time.Sleep(3 * time.Second)

	for _, v := range req {
		result, err := repo.GetStoreByName(v.Name)
		if err != nil {
			return []model.Store{}, errors.New("error get store by name after create")
		}
		res = append(res, result)
	}
	return
}

func (repo *repository) Delete(storeID int) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM stores WHERE id = $1`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.QueryContext(ctx, storeID)
	if err != nil {
		return
	}

	return
}
