package repository

import (
	"consumer-wishlist-go/model"
)

type Repositorier interface {
	Create(req []model.WishlistRequest) (err error)
}