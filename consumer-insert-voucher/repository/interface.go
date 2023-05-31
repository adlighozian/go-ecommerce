package repository

import "consumer-insert-voucher-go/model"

type Product interface {
	CreateProduct(req []model.Voucher) error
}
