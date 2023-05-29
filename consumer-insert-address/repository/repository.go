package repository

import (
	"consumer-address-go/model"
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

func (repo *repository) Create(req model.AddressRequest) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := `INSERT INTO addresses (user_id, street, city, state, country, zipcode, phone_number) values ($1, $2, $3, $4, $5, $6, $7)`
	trx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	stmt, err := trx.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.ExecContext(ctx, req.UserID, req.Street, req.City, req.State, req.Country, req.Zipcode, req.PhoneNumber)
	if err != nil {
		trx.Rollback()
		return err
	}

	trx.Commit()
	return
}
