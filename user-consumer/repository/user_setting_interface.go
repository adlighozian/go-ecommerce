package repository

import "user-consumer-go/model"

type UserSettingRepositoryI interface {
	UpdateByUserID(setting *model.UserSetting) (*model.UserSetting, error)
}
