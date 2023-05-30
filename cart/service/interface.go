//go:generate mockery --output=../mocks --name Service
package service

import (
	"cart-go/model"
)

type Servicer interface {
	Get(userID int) (res []model.Cart, err error)
	GetDetail(userID, productID int) (res model.Cart, err error)
	Create(req []model.CartRequest) (res []model.Cart, err error)
	Delete(cartID int) (err error)
}