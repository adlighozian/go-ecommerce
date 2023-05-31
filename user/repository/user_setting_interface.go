package repository

import "user-go/model"

type UserSettingRepositoryI interface {
	UpdateByUserID(newSetting *model.UserSetting) (*model.User, error)
}
