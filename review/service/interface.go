package service

import (
	"review-go/model"
)

type Servicer interface {
	GetByProductID(productID int) (res []model.Wishlist, err error)
	Create(productID int, req []model.WishlistRequest) (res []model.Wishlist, err error)
}