package model

import "time"

type PaymentMethod struct {
	Id            	 int    `json:"id"`
	PaymentGatewayID int    `json:"payment_gateway_id"`
	Name		     string `json:"name"`
}

type PaymentMethodRequest struct {
	Id            	 int    `json:"id"`
	PaymentGatewayID int    `json:"payment_gateway_id"`
	Name		     string `json:"name"`
}

type PaymentLogsRequest struct {
	UserID          int 		`json:"user_id"`
	OrderID         int 		`json:"order_id"`
	PaymentMethodID int 		`json:"payment_method_id"`
	CreatedAt       time.Time	`json:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at"`
}