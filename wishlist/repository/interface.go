package repository

import (
	"wishlist-go/model"
)

type Repositorier interface {
	Get(userID int) (res []model.Wishlist, err error)
	GetDetail(userID, wishlistID int) (res model.Wishlist, err error)
	Create(req []model.WishlistRequest) (res []model.Wishlist, err error)
	Delete(userID, wishlistID int) (err error)
}