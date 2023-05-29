package service

import "api-gateway-go/model"

type ShortenServiceI interface {
	Get(hashedURL string) (*model.APIManagement, error)
	Create(apiManagement *model.APIManagement) (*model.APIManagement, error)
}
