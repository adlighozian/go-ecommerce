package service

import (
	"cart-go/mocks"
	"cart-go/model"
	"fmt"
	"reflect"
	"testing"
)

func Test_service_Get(t *testing.T) {
	type args struct {
		userID int
	}
	tests := []struct {
		name    string
		svc     *service
		args    args
		wantRes []model.Cart
		err		error
		wantErr bool
	}{
		{
			name: "success get list of cart",
			args: args{
				userID: 1,
			},
			wantRes: []model.Cart{
				{Id: 1, UserID: 1, ProductID: 1, Quantity: 5},
				{Id: 1, UserID: 1, ProductID: 2, Quantity: 5},
			},
			err: nil,
			wantErr: false,
		},
		{
			name: "failed get list of cart",
			args: args{
				userID: 2,
			},
			wantRes: []model.Cart{},
			err: fmt.Errorf("some error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mocks.NewRepositorier(t)
			service := NewService(repoMock)

			repoMock.On("Get", tt.args.userID).Return(tt.wantRes, tt.err)

			gotRes, err := service.Get(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("service.Get() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_service_Create(t *testing.T) {
	type args struct {
		req []model.CartRequest
	}
	tests := []struct {
		name    string
		args    args
		wantRes []model.Cart
		err		error
		wantErr bool
	}{
		{
			name: "success create multiple carts",
			args: args{
				req: []model.CartRequest{
					{UserID: 1, ProductID: 1, Quantity: 5},
					{UserID: 1, ProductID: 2, Quantity: 5},
				},
			},
			wantRes: []model.Cart{
				{Id: 1, UserID: 1, ProductID: 1, Quantity: 5},
				{Id: 2, UserID: 1, ProductID: 2, Quantity: 5},
			},
			err: nil,
			wantErr: false,
		},
		{
			name: "failed create multiple carts",
			args: args{
				req: []model.CartRequest{
					{UserID: 1, ProductID: 1, Quantity: 5},
					{UserID: 1, ProductID: 2, Quantity: 5},
				},
			},
			wantRes: []model.Cart{},
			err: fmt.Errorf("some error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mocks.NewRepositorier(t)
			service := NewService(repoMock)

			repoMock.On("Create", tt.args.req).Return(tt.wantRes, tt.err)
			repoMock.On("GetDetail", tt.args.req[1].UserID, tt.args.req[1].ProductID).Return(tt.wantRes, tt.err)

			gotRes, err := service.Create(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("service.Create() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_service_Delete(t *testing.T) {
	type args struct {
		cartID int
	}
	// for mock response
	type mockGetById struct {
		res model.Cart
		err	error
	}

	tests := []struct {
		name    	string
		args    	args
		err			error
		wantErr 	bool
		getByID 	mockGetById
	}{
		{
			name: "success delete cart",
			args: args{
				cartID: 1,
			},
			err: nil,
			wantErr: false,
			// used for repo mock response
			getByID: mockGetById{
				model.Cart{Id: 1, ProductID: 1, UserID: 1, Quantity: 5}, 
				nil,
			},
		},
		{
			name: "failed delete cart",
			args: args{
				cartID: 2,
			},
			err: fmt.Errorf("some error"),
			wantErr: true,
			// used for repo mock response
			getByID: mockGetById{
				model.Cart{}, 
				fmt.Errorf("cart not found"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mocks.NewRepositorier(t)
			service := NewService(repoMock)

			repoMock.On("GetByID", tt.args.cartID).Return(tt.getByID.res, tt.getByID.err)
			repoMock.On("Delete", tt.args.cartID).Return(tt.err)

			if err := service.Delete(tt.args.cartID); (err != nil) != tt.wantErr {
				t.Errorf("service.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
