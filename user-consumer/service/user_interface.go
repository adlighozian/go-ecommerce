package service

import "user-consumer-go/model"

type UserServiceI interface {
	Create(user *model.User) (*model.User, error)
}
