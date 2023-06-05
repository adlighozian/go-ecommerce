package repository

import "size-go/model"

type Repositorier interface {
	GetSize() ([]model.Size, error)
	CreateSize(req []model.SizeReq) ([]model.Size, error)
	DeleteSize(id int) (int, error)
}
