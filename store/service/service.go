package service

import (
	"errors"
	"fmt"
	"store-go/model"
	"store-go/repository"
)

type service struct {
	repo repository.Repositorier
}

func NewService(repo repository.Repositorier) Servicer {
	return &service{
		repo: repo,
	}
}

func (svc *service) Get() (res []model.Store, err error) {
	return svc.repo.Get()
}

func (svc *service) GetStoreByName(name string) (res model.Store, err error) {
	res, err = svc.repo.GetStoreByName(name)
	if err != nil {
		return
	}

	emptyStruct := model.Store{}
	if res == emptyStruct {
		err = errors.New("store not found")
		return
	}
	return
}

func (svc *service) Create(req []model.StoreRequest) (res []model.Store, err error) {
	// check in each userID, ProductID already exist in store or not
	for _, v := range req {
		_, err := svc.GetStoreByName(v.Name)
		if err != nil {
			continue
		} else {
			err = fmt.Errorf("store with name %s already exist", v.Name)
			return []model.Store{}, err
		}
	}

	return svc.repo.Create(req)
}

func (svc *service) Delete(storeID int) (err error) {
	// check address id exist or not
	_, err = svc.repo.Get()
	if err != nil {
		return fmt.Errorf("store id %d not found", storeID)
	}
	return svc.repo.Delete(storeID)
}
