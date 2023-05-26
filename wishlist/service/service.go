package service

import (
	"errors"
	"wishlist-go/model"
	"wishlist-go/repository"
)

type service struct {
	repo repository.Repositorier
}

func NewService(repo repository.Repositorier) Servicer {
	return &service{
		repo: repo,
	}
}

func (svc *service) Get(userID int) (res []model.Wishlist, err error) {
	return svc.repo.Get(userID)
}

func (svc *service) GetDetail(userID, wishlistID int) (res model.Wishlist, err error) {
	res, err = svc.repo.GetDetail(userID, wishlistID)
	if err != nil {
		return
	}

	emptyStruct := model.Wishlist{}
	if res == emptyStruct {
		err = errors.New("wishlist not found")
		return
	}
	return
}

func (svc *service) Create(req []model.WishlistRequest) (res []model.Wishlist, err error) {
	return svc.repo.Create(req)
}

func (svc *service) Delete(userID, wishlistID int) (err error) {
	return svc.repo.Delete(userID, wishlistID)
}