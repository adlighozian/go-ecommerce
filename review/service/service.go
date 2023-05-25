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

func (svc *service) GetByProductID(productID int) (res []model.Wishlist, err error) {
	return svc.repo.GetByProductID(productID)
}

func (svc *service) Create(productID int, req []model.WishlistRequest) (res []model.Wishlist, err error) {
	return svc.repo.Create(productID, req)
}