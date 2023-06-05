package service

import "size-go/model"

type Servicer interface {
	GetSize() (model.Respon, error)
	CreateSize(req []model.SizeReq) (model.Respon, error)
	DeleteSize(idSize int) (model.Respon, error)
}
