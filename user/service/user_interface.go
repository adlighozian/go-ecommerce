package service

import "user-go/model"

type UserServiceI interface {
	GetByID(userID uint) (*model.User, error)
	UpdateByID(userID uint, profile *model.ProfileReq) (*model.User, error)
}
