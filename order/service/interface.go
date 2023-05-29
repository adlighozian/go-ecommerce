package service

import (
	"order-go/model"
)

type Servicer interface {
	GetOrders(id int) (model.Respon, error)
	CreateOrders(req model.Orders) (model.Respon, error)
	ShowOrders(idOrder int, idUser int) (model.Respon, error)
	UpdateOrders(idOrder int, req string) (model.Respon, error)
}
