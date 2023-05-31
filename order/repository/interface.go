package repository

import (
	"order-go/model"
)

type Repositorier interface {
	GetOrders(idUser int) ([]model.Orders, error)
	CreateOrders(req model.GetOrders) (model.GetOrdersSent, error)
	ShowOrders(req model.OrderItems) ([]model.OrderItem, error)
	UpdateOrders(idOrder int, req string) error
}
