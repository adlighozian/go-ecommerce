package repository

import (
	"context"
	"database/sql"
	"splash-screen-go/model"
	"splash-screen-go/publisher"
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

func (repo *repository) Get() (res []model.SplashScreen, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, image_url FROM splash_screens`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	result, err := stmt.QueryContext(ctx)
	if err != nil {
		return
	}

	for result.Next() {
		var temp model.SplashScreen
		result.Scan(&temp.Id, &temp.ImageURL)
		res = append(res, temp)
	}

	return
}

func (repo *repository) Create(req []model.SplashScreenRequest) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	query := `INSERT INTO splash_screens (image_url) values ($1)`
	trx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	stmt, err := trx.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	for _, v := range req {
		_, err = stmt.ExecContext(ctx, v.ImageURL)
		if err != nil {
			trx.Rollback()
			return
		}
	}

	trx.Commit()
	return
}

func (repo *repository) Delete(splashScreenID int) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM splash_screens WHERE id = $1`
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.QueryContext(ctx, splashScreenID)
	if err != nil {
		return
	}

	return
}
