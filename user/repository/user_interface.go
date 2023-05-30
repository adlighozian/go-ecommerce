package repository

import "user-go/model"

type UserRepositoryI interface {
	GetByID(userID uint) (*model.User, error)
	UpdateByID(profile *model.User) (*model.User, error)
}
