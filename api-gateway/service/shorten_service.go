package service

import (
	"api-gateway-go/helper/shorten"
	"api-gateway-go/model"
	"api-gateway-go/repository"
)

type ShortenService struct {
	repo    repository.ShortenRepoI
	shorten shorten.Shorten
}

func NewShortenService(repo repository.ShortenRepoI) ShortenServiceI {
	svc := new(ShortenService)
	svc.repo = repo
	svc.shorten = shorten.New()
	return svc
}

func (svc *ShortenService) Get(hashedURL string) (*model.APIManagement, error) {
	return svc.repo.Get(hashedURL)
}

func (svc *ShortenService) Create(shortenReq *model.ShortenReq) (*model.APIManagement, error) {
	url := shortenReq.EndpointURL

	url = svc.shorten.EnforceHTTP(url)
	hashedURL := svc.shorten.Encode(url)
	apiManagement := &model.APIManagement{
		APIName:           shortenReq.APIName,
		ServiceName:       shortenReq.ServiceName,
		EndpointURL:       url,
		HashedEndpointURL: hashedURL,
		IsAvailable:       shortenReq.IsAvailable,
	}
	return svc.repo.Create(apiManagement)
}
