package service

import (
	"errors"
	"net/http"
	"voucher-go/model"
	"voucher-go/repository"
)

type service struct {
	repo repository.Repositorier
}

func NewService(repo repository.Repositorier) Servicer {
	return &service{
		repo: repo,
	}
}

func (svc *service) GetVoucher() (model.Respon, error) {

	// start
	res, err := svc.repo.GetVoucher()
	if err != nil {
		return model.Respon{
			Status: http.StatusInternalServerError,
			Data:   nil,
		}, err
	}
	return model.Respon{
		Status: http.StatusOK,
		Data:   res,
	}, nil
}

func (svc *service) ShowVoucher(code string) (model.Respon, error) {

	if code == "" {
		return model.Respon{
			Status: http.StatusBadRequest,
			Data:   nil,
		}, errors.New("invalid input code voucher")
	}

	// start
	res, err := svc.repo.ShowVoucher(code)
	if err != nil {
		return model.Respon{
			Status: http.StatusInternalServerError,
			Data:   nil,
		}, err
	}
	return model.Respon{
		Status: http.StatusOK,
		Data:   res,
	}, nil
}

func (svc *service) CreateVoucher(req []model.VoucherReq) (model.Respon, error) {

	var check []model.VoucherReq

	for _, v := range req {
		if v.StoreID == 0 || v.Discount == 0 || v.Name == "" || v.StartDate == "" || v.EndDate == "" {
			continue
		}

		data := model.VoucherReq{
			StoreID:    v.StoreID,
			ProductID:  v.ProductID,
			CategoryID: v.CategoryID,
			Discount:   v.Discount,
			Name:       v.Name,
			StartDate:  v.StartDate,
			EndDate:    v.EndDate,
		}

		check = append(check, data)

	}

	if check == nil {
		return model.Respon{
			Status: http.StatusBadRequest,
			Data:   nil,
		}, errors.New("error input")
	}

	// start
	res, err := svc.repo.CreateVoucher(check)
	if err != nil {
		return model.Respon{
			Status: http.StatusInternalServerError,
			Data:   nil,
		}, err
	}
	return model.Respon{
		Status: http.StatusOK,
		Data:   res,
	}, nil
}

func (svc *service) DeleteVoucher(id int) (model.Respon, error) {
	if id <= 0 {
		return model.Respon{
			Status: http.StatusBadRequest,
			Data:   nil,
		}, errors.New("error invalid id")
	}

	// start
	res, err := svc.repo.DeleteVoucher(id)
	if err != nil {
		return model.Respon{
			Status: http.StatusInternalServerError,
			Data:   nil,
		}, err
	}
	return model.Respon{
		Status: http.StatusOK,
		Data:   res,
	}, nil
}
