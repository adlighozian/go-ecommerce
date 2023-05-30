package service

import (
	"address-go/model"
)

type Servicer interface {
	Get(userID int) (res []model.Address, err error)
	Create(req model.AddressRequest) (res []model.Address, err error)
	Delete(userID, addressID int) (err error)
}
