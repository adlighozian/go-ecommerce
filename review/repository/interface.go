package repository

import (
	"review-go/model"
)

type Repositorier interface {
	GetByProductID(productID int) (res []model.Review, err error)
	GetReviewByID(reviewID int) (res model.Review, err error)
	GetDetail(userID, productID int) (res model.Review, err error)
	Create(req []model.ReviewRequest) (res []model.Review, err error)
	Delete(reviewID int) (err error)
}
