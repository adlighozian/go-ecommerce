package service

import (
	"order-go/model"
)

type Servicer interface {
	GetOrders(idUser int) (model.Respon, error)
	GetOrdersByStoreID(storeID int) (model.Respon, error)
	CreateOrders(req model.GetOrders) (model.Respon, error)
	ShowOrders(req model.OrderItems) (model.Respon, error)
	UpdateOrders(req model.OrderUpd) (model.Respon, error)
}
