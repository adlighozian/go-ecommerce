package repository

import (
	"review-go/model"
)

type Repositorier interface {
	GetByProductID(productID int) (res []model.Wishlist, err error)
	Create(req []model.ReviewRequest) (err error)
}