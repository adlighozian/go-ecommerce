package service

import (
	"splash-screen-go/model"
)

type Servicer interface {
	Get() (res []model.SplashScreen, err error)
	Create(req []model.SplashScreenRequest) (res []model.SplashScreen, err error)
	Delete(splashScreenID int) (err error)
}
