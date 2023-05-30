package model

import "time"

type Customer struct {
	UserID int `json:"user_id"`
}

type PaymentStatusResponse struct {
	TransactionTime        string          `json:"transaction_time"`
	GrossAmount            string          `json:"gross_amount"`
	Currency               string          `json:"currency"`
	OrderID                string          `json:"order_id"`
	PaymentType            string          `json:"payment_type"`
	TransactionID          string          `json:"transaction_id"`
	TransactionStatus      string          `json:"transaction_status"`
	SettlementTime         string          `json:"settlement_time"`
	StatusMessage          string          `json:"status_message"`
	Acquirer               string          `json:"acquirer"`
	Metadata               interface{}     `json:"metadata"`
}

type PaymentLogRequest struct {
	UserID       int       `json:"user_id"`
	OrderID      int       `json:"order_id"`
	TotalPayment int64     `json:"total_payment"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
