package model

import "time"

type PaymentLogRequest struct {
	UserID          int 		`json:"user_id"`
	OrderID         int 		`json:"order_id"`
	TotalPayment	int64		`json:"total_payment"`
	PaymentMethodID int 		`json:"payment_method_id"`
	CreatedAt       time.Time	`json:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at"`
}

type PaymentLog struct {
	Id            	int    		`json:"id"`
	UserID          int 		`json:"user_id"`
	OrderID         int 		`json:"order_id"`
	TotalPayment	int64		`json:"total_payment"`
	PaymentMethodID int 		`json:"payment_method_id"`
	CreatedAt       time.Time	`json:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at"`
}