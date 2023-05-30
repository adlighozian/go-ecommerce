// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	model "cart-go/model"

	mock "github.com/stretchr/testify/mock"
)

// Repositorier is an autogenerated mock type for the Repositorier type
type Repositorier struct {
	mock.Mock
}

// Create provides a mock function with given fields: req
func (_m *Repositorier) Create(req []model.CartRequest) ([]model.Cart, error) {
	ret := _m.Called(req)

	var r0 []model.Cart
	var r1 error
	if rf, ok := ret.Get(0).(func([]model.CartRequest) ([]model.Cart, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func([]model.CartRequest) []model.Cart); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Cart)
		}
	}

	if rf, ok := ret.Get(1).(func([]model.CartRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: cartID
func (_m *Repositorier) Delete(cartID int) error {
	ret := _m.Called(cartID)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(cartID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: userID
func (_m *Repositorier) Get(userID int) ([]model.Cart, error) {
	ret := _m.Called(userID)

	var r0 []model.Cart
	var r1 error
	if rf, ok := ret.Get(0).(func(int) ([]model.Cart, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(int) []model.Cart); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Cart)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: cartID
func (_m *Repositorier) GetByID(cartID int) (model.Cart, error) {
	ret := _m.Called(cartID)

	var r0 model.Cart
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (model.Cart, error)); ok {
		return rf(cartID)
	}
	if rf, ok := ret.Get(0).(func(int) model.Cart); ok {
		r0 = rf(cartID)
	} else {
		r0 = ret.Get(0).(model.Cart)
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(cartID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDetail provides a mock function with given fields: userID, productID
func (_m *Repositorier) GetDetail(userID int, productID int) (model.Cart, error) {
	ret := _m.Called(userID, productID)

	var r0 model.Cart
	var r1 error
	if rf, ok := ret.Get(0).(func(int, int) (model.Cart, error)); ok {
		return rf(userID, productID)
	}
	if rf, ok := ret.Get(0).(func(int, int) model.Cart); ok {
		r0 = rf(userID, productID)
	} else {
		r0 = ret.Get(0).(model.Cart)
	}

	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(userID, productID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewRepositorier interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepositorier creates a new instance of Repositorier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepositorier(t mockConstructorTestingTNewRepositorier) *Repositorier {
	mock := &Repositorier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
