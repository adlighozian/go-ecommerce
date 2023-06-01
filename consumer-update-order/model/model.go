package model

type OrderUpd struct {
	OrderNumber string `json:"order_number"`
	Status      bool   `json:"status"`
}
