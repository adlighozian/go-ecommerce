package repository

import (
	"order-go/model"
)

type Repositorier interface {
	GetOrders(userID int) ([]model.Orders, error)
	ShowOrders(req model.OrderItems) (model.ResultOrders, error)
	CreateOrders(req model.GetOrders) (model.Orders, error)
	UpdateOrders(req model.OrderUpd) (model.Orders, error)
}
