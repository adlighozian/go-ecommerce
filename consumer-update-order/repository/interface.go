package repository

import "consumer-product-go/model"

type Product interface {
	UpdateProduct(req model.ProductReq) error
}
