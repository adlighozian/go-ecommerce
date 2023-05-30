package service

import "user-go/model"

type UserSettingServiceI interface {
	UpdateByUserID(userID uint, settingReq *model.SettingReq) (*model.User, error)
}
