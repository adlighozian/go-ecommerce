package model

import (
	"time"
)

type Voucher struct {
	Id         int       `json:"id"`
	StoreID    int       `json:"store_id"`
	ProductID  int       `json:"product_id"`
	CategoryID int       `json:"category_id"`
	Discount   float64   `json:"discount_value"`
	Name       string    `json:"name"`
	Code       string    `json:"code"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Created_at time.Time `json:"created_at"`
	Update_at  time.Time `json:"updated_at"`
}

type VoucherReq struct {
	// Id         int       `json:"id"`
	StoreID    int       `json:"store_id"`
	ProductID  int       `json:"product_id"`
	CategoryID int       `json:"category_id"`
	Discount   float64   `json:"discount_value"`
	Name       string    `json:"name"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
}
