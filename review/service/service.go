package service

import (
	"review-go/model"
	"review-go/repository"
)

type service struct {
	repo repository.Repositorier
}

func NewService(repo repository.Repositorier) Servicer {
	return &service{
		repo: repo,
	}
}

func (svc *service) GetByProductID(productID int) (res []model.Review, err error) {
	return svc.repo.GetByProductID(productID)
}

func (svc *service) Create(req []model.ReviewRequest) (res []model.Review, err error) {
	return svc.repo.Create(req)
}