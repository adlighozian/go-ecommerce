package model

import (
	"time"
)

type Shipping struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Created_at time.Time `json:"created_at"`
	Update_at  time.Time `json:"updated_at"`
}

type ShippingReq struct {
	Name string `json:"name"`
}
