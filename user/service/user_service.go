package service

import (
	"user-go/model"
	"user-go/repository"
)

type UserService struct {
	repo repository.UserRepositoryI
}

func NewUserService(repo repository.UserRepositoryI) UserServiceI {
	svc := new(UserService)
	svc.repo = repo
	return svc
}

func (svc *UserService) GetByID(userID uint) (*model.User, error) {
	return svc.repo.GetByID(userID)
}

func (svc *UserService) UpdateByID(userID uint, profileReq *model.ProfileReq) (*model.User, error) {
	newProfile := &model.User{
		ID:       userID,
		Username: profileReq.Username,
		FullName: profileReq.FullName,
		Age:      profileReq.Age,
		ImageURL: profileReq.ImageURL,
	}
	return svc.repo.UpdateByID(newProfile)
}
