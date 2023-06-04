package mocks

import (
	"product-go/model"

	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func NewServiceMock() *ServiceMock {
	return &ServiceMock{}
}

func (m *ServiceMock) GetProduct(req model.ProductSearch) ([]model.Product, error) {
	ret := m.Called(req)
	result := ret.Get(0).([]model.Product)
	err := ret.Error(1)
	return result, err
}
func (m *ServiceMock) ShowProduct(id int) (model.Product, error) {
	ret := m.Called(id)
	result := ret.Get(0).(model.Product)
	err := ret.Error(1)
	return result, err
}

func (m *ServiceMock) CreateProduct(req []model.Product) ([]model.Product, error) {
	ret := m.Called(req)
	result := ret.Get(0).([]model.Product)
	err := ret.Error(1)
	return result, err
}
func (m *ServiceMock) UpdateProduct(req model.ProductUpd) (model.Product, error) {
	ret := m.Called(req)
	result := ret.Get(0).(model.Product)
	err := ret.Error(1)
	return result, err
}
func (m *ServiceMock) DeleteProduct(id int) (int, error) {
	ret := m.Called(id)
	result := ret.Get(0).(int)
	err := ret.Error(1)
	return result, err
}
