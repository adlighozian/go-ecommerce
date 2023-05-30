package repository

import (
	"consumer-store-go/model"
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

func (repo *repository) Create(req []model.StoreRequest) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := `INSERT INTO stores (address_id, description, image_url, name) values ($1, $2, $3, $4)`
	trx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	stmt, err := trx.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	for _, v := range req {
		_, err := stmt.ExecContext(ctx, v.AddressID, v.Description, v.ImageURL, v.Name)
		if err != nil {
			trx.Rollback()
			return err
		}
	}

	trx.Commit()
	return
}
