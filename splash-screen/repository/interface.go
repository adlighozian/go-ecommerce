package repository

import (
	"splash-screen-go/model"
)

type Repositorier interface {
	Get() (res []model.SplashScreen, err error)
	Create(req []model.SplashScreenRequest) (err error)
	Delete(splashScreenID int) (err error)
}
