package repositoy

import (
	"order-go/model"
)

type Repositorier interface {
	GetOrders(id int) ([]model.Orders, error)
	CreateOrders(req []model.Orders) ([]model.Orders, error)
	ShowOrders(idOrder int, idUser int) (model.Orders, error)
	UpdateOrders(idOrder int, req string) error
}
