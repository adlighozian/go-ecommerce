package repository

import (
	"database/sql"
	"errors"
	"size-go/helper/failerror"
	"size-go/helper/timeout"
	"size-go/model"
	"size-go/publisher"
	"time"
)

type repository struct {
	db   *sql.DB
	sent publisher.Publisher
}

func NewRepository(db *sql.DB, sent publisher.Publisher) Repositorier {
	return &repository{
		db:   db,
		sent: sent,
	}
}

func (repo *repository) GetSize() ([]model.Size, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	query := `select * from product_sizes`

	result, err := repo.db.QueryContext(ctx, query)
	failerror.FailError(err, "fail query")

	var data = []model.Size{}
	for result.Next() {
		var temp model.Size
		result.Scan(&temp.Id, &temp.Name, &temp.Created_at, &temp.Update_at)
		data = append(data, temp)
	}

	return data, nil
}

func (repo *repository) CreateSize(sent []model.SizeReq) ([]model.Size, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	for _, v := range sent {
		var idCheck int
		queryCheck := `select id from product_sizes where name = $1`
		err := repo.db.QueryRowContext(ctx, queryCheck, v.Name).Scan(&idCheck)
		failerror.FailError(err, "error exec")

		if idCheck != 0 {
			return nil, errors.New("name " + v.Name + " already exist")
		}
	}

	err := repo.sent.Public(sent, "create_size")
	if err != nil {
		return nil, errors.New("failed publisher")
	}

	time.Sleep(1 * time.Second)

	var resultss []model.Size
	query := `select * from product_sizes where name = $1`

	stmt, err := repo.db.PrepareContext(ctx, query)
	failerror.FailError(err, "error prepare")

	for _, v := range sent {

		result, err := stmt.QueryContext(ctx, v.Name)
		failerror.FailError(err, "error prepare")

		var temp model.Size
		for result.Next() {
			result.Scan(&temp.Id, &temp.Name, &temp.Created_at, &temp.Update_at)
		}
		if temp.Id == 0 {
			continue
		}
		resultss = append(resultss, temp)
	}

	if resultss == nil {
		return nil, errors.New("error create product")
	}

	return resultss, nil
}

func (repo *repository) DeleteSize(id int) (int, error) {

	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	var idCheck int
	queryCheck := `select id from product_sizes where id = $1`
	err := repo.db.QueryRowContext(ctx, queryCheck, id).Scan(&idCheck)
	failerror.FailError(err, "error exec")

	if idCheck == 0 {
		return 0, errors.New("size not found")
	}

	query := `DELETE FROM product_sizes WHERE id = $1`
	_, err = repo.db.ExecContext(ctx, query, idCheck)
	failerror.FailError(err, "error exec")

	var idCheckDelete int
	queryCheckDelete := `select id from product_sizes where id = $1`
	err = repo.db.QueryRowContext(ctx, queryCheckDelete, idCheck).Scan(&idCheckDelete)
	failerror.FailError(err, "error exec")

	if idCheckDelete != 0 {
		return 0, errors.New("violates foreign key constraint")
	}

	return idCheck, nil
}
