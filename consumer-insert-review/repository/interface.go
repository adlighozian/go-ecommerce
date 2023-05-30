package repository

import (
	"consumer-review-go/model"
)

type Repositorier interface {
	Create(req []model.ReviewRequest) (err error)
}