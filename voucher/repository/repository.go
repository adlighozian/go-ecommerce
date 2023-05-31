package repository

import (
	"database/sql"
	"errors"
	"time"
	"voucher-go/helper/failerror"
	"voucher-go/helper/random"
	"voucher-go/helper/timeout"
	"voucher-go/model"
	"voucher-go/publisher"
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

func (repo *repository) GetVoucher(idUser int) ([]model.Voucher, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	query := `select * from voucher`

	result, err := repo.db.QueryContext(ctx, query)
	failerror.FailError(err, "fail query")

	var data []model.Voucher
	for result.Next() {
		var temp model.Voucher
		result.Scan(&temp.Id, &temp.StoreID, &temp.ProductID, &temp.CategoryID, &temp.Discount, &temp.Name, &temp.Code, &temp.StartDate, &temp.EndDate, &temp.Created_at, &temp.Update_at)
		data = append(data, temp)
	}

	return data, nil
}

func (repo *repository) ShowVoucher(code string) (model.Voucher, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	query := `select * from voucher where code = $1`

	result, err := repo.db.QueryContext(ctx, query, code)
	failerror.FailError(err, "fail query")

	var temp model.Voucher
	for result.Next() {
		result.Scan(&temp.Id, &temp.StoreID, &temp.ProductID, &temp.CategoryID, &temp.Discount, &temp.Name, &temp.Code, &temp.StartDate, &temp.EndDate, &temp.Created_at, &temp.Update_at)
	}

	if temp.Id == 0 {
		return temp, errors.New("code voucher not found")
	}

	return temp, nil
}

func (repo *repository) CreateVoucher(req []model.VoucherReq) ([]model.Voucher, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	var sent []model.Voucher

	for _, v := range req {

		startTime, err := time.Parse("2006-01-02 15:04:05", v.StartDate)
		failerror.FailError(err, "error convert")

		endTime, err := time.Parse("2006-01-02 15:04:05", v.EndDate)
		failerror.FailError(err, "error convert")

		inrandom := model.Voucher{
			StoreID:    v.StoreID,
			ProductID:  v.ProductID,
			CategoryID: v.CategoryID,
			Discount:   v.Discount,
			Name:       v.Name,
			Code:       random.NewRandom().RandomString(),
			StartDate:  startTime,
			EndDate:    endTime,
		}
		sent = append(sent, inrandom)
	}

	err := repo.sent.Public(sent, "create_voucher")
	if err != nil {
		return nil, errors.New("failed publisher")
	}

	time.Sleep(1 * time.Second)

	var resultss []model.Voucher
	query := `select * from voucher where code = $1`

	stmt, err := repo.db.PrepareContext(ctx, query)
	failerror.FailError(err, "error prepare")

	for _, v := range sent {

		result, err := stmt.QueryContext(ctx, v.Code)
		failerror.FailError(err, "error prepare")

		var temp model.Voucher
		for result.Next() {
			result.Scan(&temp.Id, &temp.StoreID, &temp.ProductID, &temp.CategoryID, &temp.Discount, &temp.Name, &temp.Code, &temp.StartDate, &temp.EndDate, &temp.Created_at, &temp.Update_at)
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

func (repo *repository) DeleteVoucher(id int) error {

	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	queryCheck := `select id from voucher where id = $1`

	result, err := repo.db.QueryContext(ctx, queryCheck, id)
	failerror.FailError(err, "fail query")

	var ids int
	for result.Next() {
		result.Scan(&ids)
	}
	if ids == 0 {
		return errors.New("id voucher not available")
	}

	query := `DELETE FROM Voucher WHERE id = $1`
	_, err = repo.db.ExecContext(ctx, query, id)
	failerror.FailError(err, "error exec")

	return nil
}
