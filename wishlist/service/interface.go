package service

import (
	"wishlist-go/model"
)

type Servicer interface {
	Get() (res []model.Wishlist, err error)
	GetDetail(id int) (res model.Wishlist, err error)
	Create(req []model.WishlistRequest) (res []model.Wishlist, err error)
	Delete(id int) (err error)
}