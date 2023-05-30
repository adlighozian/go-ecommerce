package repository

import (
	"address-go/model"
)

type Repositorier interface {
	Get(userID int) (res []model.Address, err error)
	Create(req model.AddressRequest) (res []model.Address, err error)
	Delete(addressID int) (err error)
}
