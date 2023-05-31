package service

import "user-consumer-go/model"

type UserSettingServiceI interface {
	UpdateByUserID(userSetting *model.UserSetting) (*model.UserSetting, error)
}
