package repository

import "consumer-update-order-go/model"

type Product interface {
	UpdateProduct(req model.OrderUpd) error
}
