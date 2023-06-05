package mocks

import (
	"voucher-go/model"

	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func NewServiceMock() *ServiceMock {
	return &ServiceMock{}
}

func (m *ServiceMock) GetVoucher() ([]model.Voucher, error) {
	ret := m.Called()
	result := ret.Get(0).([]model.Voucher)
	err := ret.Error(1)
	return result, err
}
func (m *ServiceMock) ShowVoucher(code string) (model.Voucher, error) {
	ret := m.Called(code)
	result := ret.Get(0).(model.Voucher)
	err := ret.Error(1)
	return result, err
}

func (m *ServiceMock) CreateVoucher(req []model.VoucherReq) ([]model.Voucher, error) {
	ret := m.Called(req)
	result := ret.Get(0).([]model.Voucher)
	err := ret.Error(1)
	return result, err
}
func (m *ServiceMock) DeleteVoucher(id int) (int, error) {
	ret := m.Called(id)
	result := ret.Get(0).(int)
	err := ret.Error(1)
	return result, err
}
