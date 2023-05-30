package service

import (
	"address-go/model"
	"address-go/repository"
	"fmt"
)

type service struct {
	repo repository.Repositorier
}

func NewService(repo repository.Repositorier) Servicer {
	return &service{
		repo: repo,
	}
}

func (svc *service) Get(userID int) (res []model.Address, err error) {
	return svc.repo.Get(userID)
}

func (svc *service) Create(req model.AddressRequest) (res []model.Address, err error) {
	return svc.repo.Create(req)
}

func (svc *service) Delete(userID, addressID int) (err error) {
	// check address id exist or not
	res, _ := svc.repo.Get(userID)
	if res != nil {
		for _, v := range res {
			if v.Id == addressID {
				return svc.repo.Delete(addressID)
			}
		}
	}
	return fmt.Errorf("address id %d in user id %d not found", addressID, userID)
}
