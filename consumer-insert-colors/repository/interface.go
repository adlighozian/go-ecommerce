package repository

import "consumer-insert-shipping-go/model"

type Product interface {
	CreateProduct(req []model.ColorsReq) error
}
