package service

import (
	"errors"
	"log"
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

	if id <= 0 {
		return model.Respon{
			Status: http.StatusBadRequest,
			Data:   nil,
		}, errors.New("invalid id product")
	}

	// start
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
		if v.StoreID == 0 || v.CategoryID == 0 || v.SizeID == 0 || v.ColorID == 0 || v.Name == "" || v.Subtitle == "" || v.Description == "" || v.UnitPrice == 0 || v.Stock == 0 || v.Weight == 0 || v.Brand == "" {
			continue
		}

		data = append(data, model.ProductReq{
			StoreID:     v.StoreID,
			CategoryID:  v.CategoryID,
			SizeID:      v.SizeID,
			ColorID:     v.ColorID,
			Name:        v.Name,
			Brand:       v.Brand,
			Subtitle:    v.Subtitle,
			Description: v.Description,
			UnitPrice:   v.UnitPrice,
			Stock:       v.Stock,
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

func (svc *service) UpdateProduct(req model.ProductUpd, id int) (model.Respon, error) {

	log.Println(id, req)

	if id <= 0 || req.Id != 0 {
		return model.Respon{
			Status: http.StatusBadRequest,
			Data:   nil,
		}, errors.New("invalid input")
	}

	data := model.ProductUpd{
		Id:          id,
		StoreID:     req.StoreID,
		CategoryID:  req.CategoryID,
		SizeID:      req.SizeID,
		ColorID:     req.ColorID,
		Name:        req.Name,
		Brand:       req.Brand,
		Subtitle:    req.Subtitle,
		Description: req.Description,
		UnitPrice:   req.UnitPrice,
		Stock:       req.Stock,
		Weight:      req.Weight,
		Status:      req.Status,
	}

	// start
	res, err := svc.repo.UpdateProduct(data)
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

func (svc *service) DeleteProduct(id int) (model.Respon, error) {
	if id <= 0 {
		return model.Respon{
			Status: http.StatusBadRequest,
			Data:   nil,
		}, errors.New("error invalid id")
	}

	// start
	res, err := svc.repo.DeleteProduct(id)
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
