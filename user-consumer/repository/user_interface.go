package repository

import "user-consumer-go/model"

type UserRepositoryI interface {
	Create(user *model.User) (*model.User, error)
	UpdateByID(user *model.User) (*model.User, error)
}
