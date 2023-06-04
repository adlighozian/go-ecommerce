package service

import (
	"fmt"
	"net/http"
	"product-go/mocks"
	"product-go/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Product_Get(t *testing.T) {
	test_get := []struct {
		name    string
		req     model.ProductSearch
		wantRes []model.Product
		err     error
		wantErr bool
	}{
		{
			name: "success get list of product",
			wantRes: []model.Product{
				{Id: 1},
			},
			err:     nil,
			wantErr: false,
		}, {
			name:    "failed get list of product",
			wantRes: []model.Product{},
			err:     fmt.Errorf("some error"),
			wantErr: true,
		},
	}
	for _, tt := range test_get {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mocks.NewServiceMock()
			service := NewService(repoMock)

			repoMock.On("GetProduct", tt.req).Return(tt.wantRes, tt.err)

			res, err := service.GetProduct(tt.req)
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
func Test_Product_Detail(t *testing.T) {
	test_get := []struct {
		name    string
		req     int
		wantRes model.Product
		err     error
		wantErr bool
	}{
		{
			name: "success get list of product",
			req:  1,
			wantRes: model.Product{
				Id: 1,
			},
			err:     nil,
			wantErr: false,
		}, {
			name:    "failed get list of product invalid id",
			wantRes: model.Product{},
			err:     fmt.Errorf("some error"),
			wantErr: true,
		}, {
			name: "success get list of product",
			req:  1,
			wantRes: model.Product{
				Id: 1,
			},
			err:     fmt.Errorf("some error"),
			wantErr: true,
		},
	}
	for _, tt := range test_get {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mocks.NewServiceMock()
			service := NewService(repoMock)

			repoMock.On("ShowProduct", tt.req).Return(tt.wantRes, tt.err)

			res, err := service.ShowProduct(tt.req)
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
func Test_Product_Create(t *testing.T) {
	test_get := []struct {
		name    string
		req     []model.ProductReq
		data    []model.Product
		wantRes []model.Product
		err     error
		wantErr bool
	}{
		// {
		// 	name: "success create product",
		// 	req: []model.ProductReq{
		// 		{
		// 			StoreID:     1,
		// 			CategoryID:  1,
		// 			ColorID:     1,
		// 			SizeID:      1,
		// 			Name:        "test",
		// 			Subtitle:    "test",
		// 			Description: "test",
		// 			UnitPrice:   1,
		// 			Stock:       1,
		// 			Weight:      1,
		// 			Brand:       "test",
		// 			ProductImage: []model.ProductImage{
		// 				{
		// 					Name:     "test",
		// 					ImageURL: "test",
		// 				},
		// 			},
		// 		},
		// 	},
		// 	wantRes: []model.Product{
		// 		{Id: 1},
		// 	},
		// 	err:     nil,
		// 	wantErr: false,
		// },
		{
			name: "failed create product invalid input",
			req: []model.ProductReq{
				{
					StoreID:     0,
					CategoryID:  1,
					ColorID:     1,
					Name:        "test",
					Subtitle:    "test",
					Description: "test",
					UnitPrice:   1,
					Stock:       1,
					Weight:      1,
					Brand:       "test",
					ProductImage: []model.ProductImage{
						{
							Name:     "test",
							ImageURL: "test",
						},
					},
				},
			},
			wantRes: []model.Product{},
			err:     fmt.Errorf("some error"),
			wantErr: true,
		},
		// {
		// 	name: "success create product",
		// 	req: []model.ProductReq{
		// 		{
		// 			StoreID:     1,
		// 			CategoryID:  1,
		// 			ColorID:     1,
		// 			Name:        "test",
		// 			Subtitle:    "test",
		// 			Description: "test",
		// 			UnitPrice:   1,
		// 			Stock:       1,
		// 			Weight:      1,
		// 			Brand:       "test",
		// 			ProductImage: []model.ProductImage{
		// 				{
		// 					Name:     "test",
		// 					ImageURL: "test",
		// 				},
		// 			},
		// 		},
		// 	},
		// 	wantRes: []model.Product{},
		// 	err:     fmt.Errorf("some error"),
		// 	wantErr: true,
		// },
	}
	for _, tt := range test_get {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mocks.NewServiceMock()
			service := NewService(repoMock)

			repoMock.On("CreateProduct", tt.req).Return(tt.wantRes, tt.err)

			res, err := service.CreateProduct(tt.req)
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

func Test_Product_Update(t *testing.T) {
	test_get := []struct {
		name    string
		req     model.ProductUpd
		id      int
		wantRes model.Product
		err     error
		wantErr bool
	}{
		{
			name: "success Update product",
			id:   1,
			req: model.ProductUpd{
				Id: 1,
				StoreID: 1,
			},
			wantRes: model.Product{
				Id: 1,
			},
			err:     nil,
			wantErr: false,
		},
		{
			name:    "failed Update product invalid id",
			wantRes: model.Product{},
			err:     fmt.Errorf("some error"),
			wantErr: true,
		}, {
			name: "success Update product",
			id:   1,
			req: model.ProductUpd{
				Id: 1,
				StoreID: 1,
			},
			wantRes: model.Product{
				Id: 1,
			},
			err:     fmt.Errorf("some error"),
			wantErr: true,
		},
	}
	for _, tt := range test_get {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mocks.NewServiceMock()
			service := NewService(repoMock)

			repoMock.On("UpdateProduct", tt.req).Return(tt.wantRes, tt.err)

			res, err := service.UpdateProduct(tt.req, tt.id)
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
