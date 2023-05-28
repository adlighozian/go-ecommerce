//go:generate mockery --output=../mocks --name Repositorier
package repository

import (
	"wishlist-go/model"
)

type Repositorier interface {
	Get(userID int) (res []model.Wishlist, err error)
	GetByID(wishlistID int) (res model.Wishlist, err error)
	GetDetail(userID, productID int) (res model.Wishlist, err error)
	Create(req []model.WishlistRequest) (res []model.Wishlist, err error)
	Delete(wishlistID int) (err error)
}