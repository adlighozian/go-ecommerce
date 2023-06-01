package service

import (
	"fmt"
	"review-go/model"
	"review-go/repository"
)

type service struct {
	repo repository.Repositorier
}

func NewService(repo repository.Repositorier) Servicer {
	return &service{
		repo: repo,
	}
}

func (svc *service) GetByProductID(productID int) (res []model.Review, err error) {
	return svc.repo.GetByProductID(productID)
}

func (svc *service) Create(req []model.ReviewRequest) (res []model.Review, err error) {
	return svc.repo.Create(req)
}

func (svc *service) Delete(reviewID int) (err error) {
	// check cart id exist or not
	emptyStruct := model.Review{}
	res, _ := svc.repo.GetReviewByID(reviewID)
	if res == emptyStruct {
		return fmt.Errorf("item with id %d not found", reviewID)
	}
	return svc.repo.Delete(reviewID)
}
