package service

import (
	"user-consumer-go/model"
	"user-consumer-go/repository"
)

type UserSettingService struct {
	repo repository.UserSettingRepositoryI
}

func NewUserSettingService(repo repository.UserSettingRepositoryI) UserSettingServiceI {
	svc := new(UserSettingService)
	svc.repo = repo
	return svc
}

func (svc *UserSettingService) UpdateByUserID(userSetting *model.UserSetting) (*model.UserSetting, error) {
	return svc.repo.UpdateByUserID(userSetting)
}
