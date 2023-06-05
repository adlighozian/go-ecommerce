package repository

import "product-colors-go/model"

type Repositorier interface {
	GetColors() ([]model.Colors, error)
	CreateColors(req []model.ColorsReq) ([]model.Colors, error)
	DeleteColors(id int) (int, error)
}
