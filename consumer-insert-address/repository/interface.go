package repository

import (
	"consumer-address-go/model"
)

type Repositorier interface {
	Create(req model.AddressRequest) (err error)
}
