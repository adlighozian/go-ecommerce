package repositoy

import (
	"review-go/model"
)

type Repositorier interface {
	Get() (res []model.Review, err error)
	GetDetail(id int) (res model.Review, err error)
	Create(req []model.ReviewRequest) (err error)
}