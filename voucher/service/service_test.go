package service

import (
	"fmt"
	"net/http"
	"testing"
	"voucher-go/mocks"
	"voucher-go/model"

	"github.com/stretchr/testify/require"
)

func TestServiceGet(t *testing.T) {
	test_get := []struct {
		name    string
		svc     *service
		wantRes []model.Voucher
		err     error
		wantErr bool
	}{
		{
			name: "success get list of voucher",
			wantRes: []model.Voucher{
				{Id: 1},
				{Id: 1},
			},
			err:     nil,
			wantErr: false,
		}, {
			name:    "failed get list of voucher",
			wantRes: []model.Voucher{},
			err:     fmt.Errorf("some error"),
			wantErr: true,
		},
	}
	for _, tt := range test_get {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mocks.NewServiceMock()
			service := NewService(repoMock)

			repoMock.On("GetVoucher").Return(tt.wantRes, tt.err)

			res, err := service.GetVoucher()
			if !tt.wantErr {
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, res.Status)
			} else {
				require.Error(t, err)
				require.Equal(t, http.StatusInternalServerError, res.Status)
			}
		})
	}
}

func TestServiceDetail(t *testing.T) {
	test_get := []struct {
		name    string
		svc     *service
		code    string
		wantRes model.Voucher
		err     error
		wantErr bool
	}{
		{
			name: "success get detail voucher",
			code: "9HUrhC0jiO",
			wantRes: model.Voucher{
				Id: 1,
			},
			err:     nil,
			wantErr: false,
		}, {
			name:    "failed get detail voucher no code",
			code:    "",
			wantRes: model.Voucher{},
			err:     fmt.Errorf("some error"),
			wantErr: true,
		}, {
			name:    "failed get detail voucher no code",
			code:    "ajsdfhkuj",
			wantRes: model.Voucher{},
			err:     fmt.Errorf("some error"),
			wantErr: true,
		},
	}
	for _, tt := range test_get {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mocks.NewServiceMock()
			service := NewService(repoMock)

			repoMock.On("ShowVoucher", tt.code).Return(tt.wantRes, tt.err)

			res, err := service.ShowVoucher(tt.code)
			if res.Status == http.StatusOK {
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, res.Status)
			} else if res.Status == http.StatusBadRequest {
				require.Error(t, err)
				require.Equal(t, http.StatusBadRequest, res.Status)
			} else if res.Status == http.StatusInternalServerError {
				require.Error(t, err)
				require.Equal(t, http.StatusInternalServerError, res.Status)
			}
		})
	}
}
func TestServiceCreate(t *testing.T) {
	test_get := []struct {
		name    string
		svc     *service
		req     []model.VoucherReq
		wantRes []model.Voucher
		err     error
		wantErr bool
	}{
		{
			name: "success create voucher",
			req: []model.VoucherReq{
				{
					StoreID:   1,
					Discount:  1,
					Name:      "test",
					StartDate: "01",
					EndDate:   "01",
				},
			},
			wantRes: []model.Voucher{
				{
					Id: 1,
				},
			},
			err:     nil,
			wantErr: false,
		}, {
			name: "failed input create voucher",
			req: []model.VoucherReq{
				{
					StoreID: 0,
				},
			},
			wantRes: []model.Voucher{},
			err:     fmt.Errorf("some error"),
			wantErr: true,
		}, {
			name: "failed get list of voucher no code",
			req: []model.VoucherReq{
				{
					StoreID:   1,
					Discount:  1,
					Name:      "test",
					StartDate: "01",
					EndDate:   "01",
				},
			},
			wantRes: []model.Voucher{},
			err:     fmt.Errorf("some error"),
			wantErr: true,
		},
	}
	for _, tt := range test_get {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mocks.NewServiceMock()
			service := NewService(repoMock)

			repoMock.On("CreateVoucher", tt.req).Return(tt.wantRes, tt.err)

			res, err := service.CreateVoucher(tt.req)
			if res.Status == http.StatusOK {
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, res.Status)
			} else if res.Status == http.StatusBadRequest {
				require.Error(t, err)
				require.Equal(t, http.StatusBadRequest, res.Status)
			} else if res.Status == http.StatusInternalServerError {
				require.Error(t, err)
				require.Equal(t, http.StatusInternalServerError, res.Status)
			}
		})
	}
}
func TestServiceDelete(t *testing.T) {
	test_get := []struct {
		name    string
		svc     *service
		req     int
		wantRes int
		err     error
		wantErr bool
	}{
		{
			name:    "success delete voucher",
			req:     1,
			wantRes: 1,
			err:     nil,
			wantErr: false,
		}, {
			name:    "failed delete voucher no id",
			req:     0,
			wantRes: 0,
			err:     fmt.Errorf("some error"),
			wantErr: true,
		}, {
			name:    "failed get list of voucher no code",
			req:     1,
			wantRes: 1,
			err:     fmt.Errorf("some error"),
			wantErr: true,
		},
	}
	for _, tt := range test_get {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mocks.NewServiceMock()
			service := NewService(repoMock)

			repoMock.On("DeleteVoucher", tt.req).Return(tt.wantRes, tt.err)

			res, err := service.DeleteVoucher(tt.req)
			if res.Status == http.StatusOK {
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, res.Status)
			} else if res.Status == http.StatusBadRequest {
				require.Error(t, err)
				require.Equal(t, http.StatusBadRequest, res.Status)
			} else if res.Status == http.StatusInternalServerError {
				require.Error(t, err)
				require.Equal(t, http.StatusInternalServerError, res.Status)
			}
		})
	}
}
