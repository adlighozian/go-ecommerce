package service

import (
	"review-go/model"
)

type Servicer interface {
	GetByProductID(productID int) (res []model.Review, err error)
	Create(req []model.ReviewRequest) (res []model.Review, err error)
	Delete(reviewID int) (err error)
}