//go:generate mockery --output=../mocks --name ShortenServiceI
package service

import "api-gateway-go/model"

type ShortenServiceI interface {
	Get(hashedURL string) (*model.APIManagement, error)
	Create(shortenReq *model.ShortenReq) (*model.APIManagement, error)
}
