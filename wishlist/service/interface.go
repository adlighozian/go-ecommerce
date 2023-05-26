//go:generate mockery --output=../mocks --name Servicer
package service

import (
	"wishlist-go/model"
)

type Servicer interface {
	Get(userID int) (res []model.Wishlist, err error)
	GetDetail(userID, wishlistID int) (res model.Wishlist, err error)
	Create(req []model.WishlistRequest) (res []model.Wishlist, err error)
	Delete(userID, wishlistID int) (err error)
}