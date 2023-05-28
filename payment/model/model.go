package model

import "time"

type PaymentMethod struct {
	Id            	 int    `json:"id"`
	PaymentGatewayID string `json:"payment_gateway_id"`
	Name		     string `json:"name"`
}

type PaymentMethodRequest struct {
	PaymentGatewayID string `json:"payment_gateway_id"`
	Name		     string `json:"name"`
}

type PaymentLogRequest struct {
	UserID          int 		`json:"user_id"`
	OrderID         int 		`json:"order_id"`
	PaymentMethodID int 		`json:"payment_method_id"`
	CreatedAt       time.Time	`json:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at"`
}

type PaymentLog struct {
	Id            	int    `json:"id"`
	UserID          int 		`json:"user_id"`
	OrderID         int 		`json:"order_id"`
	PaymentMethodID int 		`json:"payment_method_id"`
	CreatedAt       time.Time	`json:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at"`
}