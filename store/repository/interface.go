package repository

import (
	"store-go/model"
)

type Repositorier interface {
	Get() (res []model.Store, err error)
	GetStoreByName(name string) (res model.Store, err error)
	Create(req []model.StoreRequest) (res []model.Store, err error)
	Delete(storeID int) (err error)
}
