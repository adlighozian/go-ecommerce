package model

import (
	"time"
)

type Product struct {
	Id           int            `json:"id"`
	StoreID      int            `json:"store_id"`
	CategoryID   int            `json:"category_id"`
	SizeID       int            `json:"size_id"`
	ColorID      int            `json:"color_id"`
	Name         string         `json:"name"`
	Brand        string         `json:"brand"`
	Subtitle     string         `json:"subtitle"`
	Description  string         `json:"description"`
	UnitPrice    float64        `json:"unit_price"`
	Status       bool           `json:"status"`
	Stock        int            `json:"stock"`
	Sku          string         `json:"sku"`
	Weight       float64        `json:"weight"`
	ProductImage []ProductImage `json:"image"`
	Created_at   time.Time      `json:"created_at"`
	Update_at    time.Time      `json:"updated_at"`
}

type ProductImage struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

