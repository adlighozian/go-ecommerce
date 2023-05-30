package service

import (
	"user-consumer-go/model"
	"user-consumer-go/repository"
)

type UserService struct {
	repo repository.UserRepositoryI
}

func NewUserService(repo repository.UserRepositoryI) UserServiceI {
	svc := new(UserService)
	svc.repo = repo
	return svc
}

func (svc *UserService) Create(user *model.User) (*model.User, error) {
	return svc.repo.Create(user)
}

func (svc *UserService) UpdateByID(user *model.User) (*model.User, error) {
	return svc.repo.UpdateByID(user)
}
