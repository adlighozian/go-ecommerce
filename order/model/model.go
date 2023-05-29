package model

import "time"

type Orders struct {
	Id              int       `json:"id"`
	UserID          int       `json:"user_id"`
	PaymentMethodID int       `json:"payment_method_id"`
	ShippingID      int       `json:"shipping_id"`
	TotalPrice      float64   `json:"total_price"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type OrderItems struct {
	Id int
	OrderId int	
}

type Respon struct {
	Status int
	Data   any
}

type ResponSuccess struct {
	Status  int
	Message string
	Data    any
}
type ResponError struct {
	Status  int
	Message string
	Error   string
}


