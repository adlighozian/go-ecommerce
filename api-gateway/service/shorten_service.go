package service

import (
	"api-gateway-go/model"
	"api-gateway-go/repository"
)

type ShortenService struct {
	repo repository.ShortenRepoI
}

func NewShortenService(repo repository.ShortenRepoI) ShortenServiceI {
	svc := new(ShortenService)
	svc.repo = repo
	return svc
}

func (svc *ShortenService) Get(hashedURL string) (*model.APIManagement, error) {
	return svc.repo.Get(hashedURL)
}

func (svc *ShortenService) Create(apiManagement *model.APIManagement) (*model.APIManagement, error) {
	return svc.repo.Create(apiManagement)
}
