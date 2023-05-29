package repository

import "api-gateway-go/model"

type ShortenRepoI interface {
	Get(hashedURL string) (*model.APIManagement, error)
	Create(apiManagement *model.APIManagement) (*model.APIManagement, error)
}
