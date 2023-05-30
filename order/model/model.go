package model

import "time"

type Orders struct {
	Id          int       `json:"id"`
	UserID      int       `json:"user_id"`
	ShippingID  int       `json:"shipping_id"`
	TotalPrice  float64   `json:"total_price"`
	Status      bool      `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	OrderNumber string    `json:"order_number"`
}

type OrderReq struct {
	UserID     int
	ShippingID int
	TotalPrice float64
}

type OrderItems struct {
	UserId      int
	OrderNumber string
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
