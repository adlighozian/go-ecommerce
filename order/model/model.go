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

type OrderItem struct {
	Id         int       `json:"id"`
	OrderID    int       `json:"order_id"`
	ProductId  int       `json:"product_id"`
	Quantity   int       `json:"quantity"`
	TotalPrice float64   `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type GetOrders struct {
	UserID       int            `json:"user_id"`
	ShippingID   int            `json:"shipping_id"`
	TotalPrice   float64        `json:"total_price"`
	OrderItemReq []OrderItemReq `json:"order_items"`
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
