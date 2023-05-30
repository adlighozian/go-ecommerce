package service

import (
	"errors"
	"fmt"
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

func (svc *service) GetDetail(userID, productID int) (res model.Wishlist, err error) {
	res, err = svc.repo.GetDetail(userID, productID)
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
	// check in each userID, ProductID already exist in wishlist or not
	for _, v := range req {
		_, err := svc.GetDetail(v.UserID, v.ProductID)
		if err != nil {
			continue
		} else {
			err = fmt.Errorf("wishlist with product_id %d in user_id %d already exist", v.ProductID, v.UserID)
			return []model.Wishlist{}, err
		}
	}

	return svc.repo.Create(req)
}

func (svc *service) Delete(wishlistID int) (err error) {
	// check wishlist id exist or not
	emptyStruct := model.Wishlist{}
	res, _ := svc.repo.GetByID(wishlistID)
	if res == emptyStruct {
		return fmt.Errorf("item with id %d not found", wishlistID)
	}
	
	return svc.repo.Delete(wishlistID)
}