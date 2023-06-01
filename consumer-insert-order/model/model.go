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

type OrderItemReq struct {
	ProductId  int     `json:"product_id"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
}

type GetOrdersSent struct {
	UserID       int            `json:"user_id"`
	ShippingID   int            `json:"shipping_id"`
	TotalPrice   float64        `json:"total_price"`
	Status       bool           `json:"status"`
	OrderNumber  string         `json:"order_number"`
	OrderItemReq []OrderItemReq `json:"order_items"`
}
