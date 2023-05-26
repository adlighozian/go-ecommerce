package model

import (
	"time"
)

type Product struct {
	Id          int
	StoreID     int
	CategoryID  int
	SizeID      int
	ColorID     int
	Name        string
	Subtitle    string
	Description string
	UnitPrice   float64
	Status      bool
	Stock       int
	Sku         string
	Weight      float64
	Created_at  time.Time
	Update_at   time.Time
}

type ProductReq struct {
	Id          int
	StoreID     int
	CategoryID  int
	SizeID      int
	ColorID     int
	Name        string
	Subtitle    string
	Description string
	UnitPrice   float64
	Status      *bool
	Stock       int
	Sku         string
	Weight      float64
}

type Respon struct {
	Status int
	Data   any
}

type ProductSearch struct {
	Arraival string
	Brand    string
	Category string
	Name     string
}
