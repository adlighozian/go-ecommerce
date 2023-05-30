package repository

import (
	"consumer-store-go/model"
)

type Repositorier interface {
	Create(req []model.StoreRequest) (err error)
}
