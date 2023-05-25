package repository

import (
	"wishlist-go/model"
)

type Repositorier interface {
	Get() (res []model.Wishlist, err error)
	GetDetail(id int) (res model.Wishlist, err error)
	Create(req []model.WishlistRequest) (res []model.Wishlist, err error)
	Delete(id int) (err error)
}