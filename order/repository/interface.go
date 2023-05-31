package repository

import (
	"order-go/model"
)

type Repositorier interface {
	GetOrders(idUser int) ([]model.Orders, error)
	CreateOrders(req []model.OrderReq) ([]model.Orders, error)
	ShowOrders(req model.OrderItems) (model.Orders, error)
	UpdateOrders(idOrder int, req string) error
}
