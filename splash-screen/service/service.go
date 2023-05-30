package service

import (
	"fmt"
	"splash-screen-go/model"
	"splash-screen-go/repository"
)

type service struct {
	repo repository.Repositorier
}

func NewService(repo repository.Repositorier) Servicer {
	return &service{
		repo: repo,
	}
}

func (svc *service) Get() (res []model.SplashScreen, err error) {
	return svc.repo.Get()
}

func (svc *service) Create(req []model.SplashScreenRequest) (res []model.SplashScreen, err error) {
	err = svc.repo.Create(req)
	if err != nil {
		return nil, fmt.Errorf("error create : %s", err.Error())
	}
	return svc.repo.Get()
}

func (svc *service) Delete(splashScreenID int) (err error) {
	// check address id exist or not
	_, err = svc.repo.Get()
	if err != nil {
		return fmt.Errorf("splash screen id %d not found", splashScreenID)
	}
	return svc.repo.Delete(splashScreenID)
}
