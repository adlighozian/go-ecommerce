package service

import (
	"user-go/model"
	"user-go/repository"
)

type UserSettingService struct {
	repo repository.UserSettingRepositoryI
}

func NewUserSettingService(repo repository.UserSettingRepositoryI) UserSettingServiceI {
	svc := new(UserSettingService)
	svc.repo = repo
	return svc
}

func (svc *UserSettingService) UpdateByUserID(userID uint, settingReq *model.SettingReq) (*model.User, error) {
	newSetting := &model.UserSetting{
		UserID:       userID,
		Notification: settingReq.Notification,
		DarkMode:     settingReq.DarkMode,
		LanguageID:   settingReq.LanguageID,
	}
	return svc.repo.UpdateByUserID(newSetting)
}
