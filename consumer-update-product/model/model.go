package model

type ProductUpd struct {
	Id          int     `json:"id"`
	StoreID     int     `json:"store_id"`
	CategoryID  int     `json:"category_id"`
	SizeID      int     `json:"size_id"`
	ColorID     int     `json:"color_id"`
	Name        string  `json:"name"`
	Brand       string  `json:"brand"`
	Subtitle    string  `json:"subtitle"`
	Description string  `json:"description"`
	UnitPrice   float64 `json:"unit_price"`
	Status      bool    `json:"status"`
	Stock       int     `json:"stock"`
	Weight      float64 `json:"weight"`
}
