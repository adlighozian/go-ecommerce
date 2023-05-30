package service

import (
	"errors"
	"net/http"
	"product-go/model"
	"product-go/repository"
)

type service struct {
	repo repository.Repositorier
}

func NewService(repo repository.Repositorier) Servicer {
	return &service{
		repo: repo,
	}
}

func (svc *service) GetProduct(req model.ProductSearch) (model.Respon, error) {
	res, err := svc.repo.GetProduct(req)
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

func (svc *service) ShowProduct(id int) (model.Respon, error) {

	if id == 0 {
		return model.Respon{
			Status: http.StatusBadRequest,
			Data:   nil,
		}, errors.New("error detail, id not found")
	}

	res, err := svc.repo.ShowProduct(id)
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

func (svc *service) CreateProduct(req []model.ProductReq) (model.Respon, error) {

	var data []model.ProductReq

	for _, v := range req {
		if v.StoreID == 0 || v.CategoryID == 0 || v.SizeID == 0 || v.ColorID == 0 || v.Name == "" || v.Subtitle == "" || v.Description == "" || v.UnitPrice == 0 || v.Stock == 0 || v.Weight == 0 {
			continue
		}
		data = append(data, model.ProductReq{
			StoreID:     v.StoreID,
			CategoryID:  v.CategoryID,
			SizeID:      v.SizeID,
			ColorID:     v.ColorID,
			Name:        v.Name,
			Subtitle:    v.Subtitle,
			Description: v.Description,
			UnitPrice:   v.UnitPrice,
			Status:      v.Status,
			Stock:       v.Stock,
			Sku:         v.Sku,
			Weight:      v.Weight,
		})
	}

	if data == nil {
		return model.Respon{
			Status: http.StatusBadRequest,
			Data:   nil,
		}, errors.New("error input")
	}

	// start
	res, err := svc.repo.CreateProduct(data)
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

func (svc *service) UpdateProduct(req model.ProductReq) (model.Respon, error) {

	if req.Id == 0 {
		return model.Respon{
			Status: http.StatusBadRequest,
			Data:   nil,
		}, errors.New("error update data")
	}

	// start
	err := svc.repo.UpdateProduct(req)
	if err != nil {
		return model.Respon{
			Status: http.StatusInternalServerError,
			Data:   nil,
		}, err
	}
	return model.Respon{
		Status: http.StatusOK,
		Data:   nil,
	}, nil
}

func (svc *service) DeleteProduct(id int) (model.Respon, error) {
	if id == 0 {
		return model.Respon{
			Status: http.StatusBadRequest,
			Data:   nil,
		}, errors.New("error invalid id")
	}

	err := svc.repo.DeleteProduct(id)
	if err != nil {
		return model.Respon{
			Status: http.StatusInternalServerError,
			Data:   nil,
		}, err
	}
	return model.Respon{
		Status: http.StatusOK,
		Data:   nil,
	}, nil
}
